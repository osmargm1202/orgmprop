package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orgmai/internal/chat"
	"orgmai/internal/config"
	"orgmai/internal/conversation"
	"orgmai/internal/logger"
	"orgmai/internal/ui"
)

var (
	debugFlag = false
)

func main() {
	// Inicializar logger
	logDir := filepath.Join(config.ConfigDir, "logs")
	if err := logger.Init(logDir, debugFlag); err != nil {
		fmt.Fprintf(os.Stderr, "Error inicializando logger: %v\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]

	// Verificar flag --debug
	for i, arg := range args {
		if arg == "--debug" {
			debugFlag = true
			logger.Init(logDir, true) // Reinicializar con debug
			args = append(args[:i], args[i+1:]...)
			break
		}
	}

	if len(args) == 0 {
		showHelp()
		return
	}

	command := args[0]
	remainingArgs := args[1:]

	logger.Debug("Comando ejecutado: %s, args: %v", command, remainingArgs)

	switch command {
	case "apikey":
		cmdApikey()
	case "config":
		cmdConfig()
	case "chat":
		cmdChat(remainingArgs)
	case "prev":
		cmdPrev()
	default:
		// Si no es un comando reconocido, tratar como pregunta directa
		cmdChat(args)
	}
}

func showHelp() {
	fmt.Println("orgmai - CLI para interactuar con Claude")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  orgmai apikey              - Configurar API key de Claude")
	fmt.Println("  orgmai config              - Seleccionar modelo Claude")
	fmt.Println("  orgmai chat [pregunta]     - Iniciar/continuar conversación")
	fmt.Println("  orgmai prev                - Seleccionar conversación anterior")
	fmt.Println()
	fmt.Println("Opciones:")
	fmt.Println("  --debug                    - Mostrar logs de debug en consola")
}

func cmdApikey() {
	ui.PrintInfo("Configura tu API key de Claude\n")

	apiKey, err := ui.Input("Ingrese su Claude API Key", "sk-ant-...")
	if err != nil {
		ui.PrintError("Error obteniendo input: " + err.Error())
		logger.Error("Error obteniendo API key: %v", err)
		os.Exit(1)
	}

	if apiKey == "" {
		ui.PrintWarning("Operación cancelada")
		return
	}

	cfg, err := config.Load()
	if err != nil {
		ui.PrintError("Error cargando configuración: " + err.Error())
		logger.Error("Error cargando configuración: %v", err)
		os.Exit(1)
	}

	cfg.ClaudeAPIKey = apiKey
	if err := config.Save(cfg); err != nil {
		ui.PrintError("Error guardando configuración: " + err.Error())
		logger.Error("Error guardando configuración: %v", err)
		os.Exit(1)
	}

	ui.PrintSuccess("API key guardada correctamente")
	logger.Info("API key configurada")
}

func cmdConfig() {
	ui.PrintInfo("Selecciona el modelo de Claude:\n")

	cfg, err := config.Load()
	if err != nil {
		ui.PrintError("Error cargando configuración: " + err.Error())
		logger.Error("Error cargando configuración: %v", err)
		os.Exit(1)
	}

	currentModel := cfg.Model
	if currentModel == "" {
		currentModel = config.DefaultModel
	}

	models := config.AvailableModels()
	options := make([]string, len(models))
	for i, model := range models {
		if model == currentModel {
			options[i] = model + " (actual)"
		} else {
			options[i] = model
		}
	}

	selected, err := ui.Select("Modelo:", options)
	if err != nil {
		ui.PrintError("Error en selección: " + err.Error())
		logger.Error("Error en selección: %v", err)
		os.Exit(1)
	}

	if selected == "" {
		return
	}

	// Remover "(actual)" si está presente
	modelName := strings.TrimSuffix(selected, " (actual)")

	cfg.Model = modelName
	if err := config.Save(cfg); err != nil {
		ui.PrintError("Error guardando configuración: " + err.Error())
		logger.Error("Error guardando configuración: %v", err)
		os.Exit(1)
	}

	ui.PrintSuccess(fmt.Sprintf("Configuración guardada: %s", modelName))
	logger.Info("Modelo configurado: %s", modelName)
}

func cmdChat(args []string) {
	// Obtener API key
	apiKey, err := config.GetAPIKey()
	if err != nil {
		ui.PrintError(err.Error())
		logger.Error("API key no configurada: %v", err)
		os.Exit(1)
	}

	// Obtener modelo
	model, err := config.GetModel()
	if err != nil {
		ui.PrintError("Error obteniendo modelo: " + err.Error())
		logger.Error("Error obteniendo modelo: %v", err)
		os.Exit(1)
	}

	// Crear cliente
	client := chat.NewClient(apiKey, model)

	// Si hay argumentos, es una pregunta directa
	if len(args) > 0 {
		pregunta := strings.Join(args, " ")
		handleDirectQuestion(client, model, pregunta)
		return
	}

	// Modo interactivo
	handleInteractiveChat(client, model)
}

func handleDirectQuestion(client *chat.Client, model, pregunta string) {
	ui.PrintInfo(fmt.Sprintf("Modelo: %s\n", model))
	ui.PrintInfo("Respuesta:\n")

	messages := []conversation.Message{
		{Role: "user", Content: pregunta},
	}

	ctx := context.Background()
	response, err := client.SendStream(ctx, messages)
	if err != nil {
		ui.PrintError("Error en la solicitud: " + err.Error())
		logger.Error("Error en solicitud: %v", err)
		os.Exit(1)
	}

	if response != "" {
		messages = append(messages, conversation.Message{
			Role:    "assistant",
			Content: response,
		})

		// Guardar conversación
		now := time.Now()
		timestamp := now.Format("2006-01-02-15-04-05")
		title := conversation.GenerateTitle(messages)
		if err := conversation.Save(messages, title, timestamp); err != nil {
			ui.PrintWarning("Error guardando conversación: " + err.Error())
			logger.Warn("Error guardando conversación: %v", err)
		} else {
			ui.PrintSuccess(fmt.Sprintf("Conversación guardada: %s-%s.md", timestamp, title))
		}

		// Continuar en modo interactivo
		handleInteractiveChatWithMessages(client, model, messages, title, timestamp)
	}
}

func handleInteractiveChat(client *chat.Client, model string) {
	var messages []conversation.Message
	var title string
	var timestamp string

	ui.PrintInfo(fmt.Sprintf("Modelo: %s\n", model))

	for {
		userInput, err := ui.Input("Tú>", "")
		if err != nil {
			logger.Debug("Input cancelado o error: %v", err)
			break
		}

		if userInput == "" {
			break
		}

		if strings.ToLower(userInput) == "salir" || strings.ToLower(userInput) == "exit" || strings.ToLower(userInput) == "quit" {
			break
		}

		messages = append(messages, conversation.Message{
			Role:    "user",
			Content: userInput,
		})

		ui.PrintInfo("\nRespuesta:\n")

		ctx := context.Background()
		response, err := client.SendStream(ctx, messages)
		if err != nil {
			ui.PrintError("Error en la solicitud: " + err.Error())
			logger.Error("Error en solicitud: %v", err)
			continue
		}

		if response != "" {
			messages = append(messages, conversation.Message{
				Role:    "assistant",
				Content: response,
			})

			// Guardar conversación
			if len(messages) == 2 {
				// Primera respuesta, generar título y timestamp
				now := time.Now()
				timestamp = now.Format("2006-01-02-15-04-05")
				title = conversation.GenerateTitle(messages)
			}

			if timestamp != "" {
				if err := conversation.Save(messages, title, timestamp); err != nil {
					ui.PrintWarning("Error guardando conversación: " + err.Error())
					logger.Warn("Error guardando conversación: %v", err)
				}
			}
		}
	}
}

func handleInteractiveChatWithMessages(client *chat.Client, model string, initialMessages []conversation.Message, title, timestamp string) {
	messages := initialMessages

	for {
		userInput, err := ui.Input("Tú>", "")
		if err != nil {
			logger.Debug("Input cancelado o error: %v", err)
			break
		}

		if userInput == "" {
			break
		}

		if strings.ToLower(userInput) == "salir" || strings.ToLower(userInput) == "exit" || strings.ToLower(userInput) == "quit" {
			break
		}

		messages = append(messages, conversation.Message{
			Role:    "user",
			Content: userInput,
		})

		ui.PrintInfo("\nRespuesta:\n")

		ctx := context.Background()
		response, err := client.SendStream(ctx, messages)
		if err != nil {
			ui.PrintError("Error en la solicitud: " + err.Error())
			logger.Error("Error en solicitud: %v", err)
			continue
		}

		if response != "" {
			messages = append(messages, conversation.Message{
				Role:    "assistant",
				Content: response,
			})

			// Guardar conversación
			if err := conversation.Save(messages, title, timestamp); err != nil {
				ui.PrintWarning("Error guardando conversación: " + err.Error())
				logger.Warn("Error guardando conversación: %v", err)
			}
		}
	}
}

func cmdPrev() {
	conversations, err := conversation.List()
	if err != nil {
		ui.PrintError("Error listando conversaciones: " + err.Error())
		logger.Error("Error listando conversaciones: %v", err)
		os.Exit(1)
	}

	if len(conversations) == 0 {
		ui.PrintWarning("No hay conversaciones anteriores")
		return
	}

	ui.PrintInfo("Selecciona una conversación para continuar:\n")

	// Formatear opciones
	options := make([]string, len(conversations))
	for i, conv := range conversations {
		// Parsear nombre de archivo: timestamp-title.md
		parts := strings.Split(strings.TrimSuffix(conv, ".md"), "-")
		if len(parts) >= 5 {
			title := strings.Join(parts[4:], "-")
			dateStr := fmt.Sprintf("%s-%s-%s %s", parts[0], parts[1], parts[2], parts[3])
			options[i] = fmt.Sprintf("%s (%s)", title, dateStr)
		} else {
			options[i] = conv
		}
	}

	selected, err := ui.Select("Conversación:", options)
	if err != nil {
		ui.PrintError("Error en selección: " + err.Error())
		logger.Error("Error en selección: %v", err)
		os.Exit(1)
	}

	if selected == "" {
		return
	}

	// Encontrar el índice seleccionado
	idx := -1
	for i, opt := range options {
		if opt == selected {
			idx = i
			break
		}
	}

	if idx == -1 {
		ui.PrintError("Conversación no encontrada")
		return
	}

	filename := conversations[idx]

	// Cargar conversación
	messages, err := conversation.Load(filename)
	if err != nil {
		ui.PrintError("Error cargando conversación: " + err.Error())
		logger.Error("Error cargando conversación: %v", err)
		os.Exit(1)
	}

	// Parsear timestamp y título
	parts := strings.Split(strings.TrimSuffix(filename, ".md"), "-")
	var timestamp, title string
	if len(parts) >= 5 {
		timestamp = strings.Join(parts[:4], "-")
		title = strings.Join(parts[4:], "-")
	} else {
		timestamp = time.Now().Format("2006-01-02-15-04-05")
		title = "conversacion"
	}

	// Mostrar conversación
	filepath := filepath.Join(conversation.ConversationsDir, filename)
	data, err := os.ReadFile(filepath)
	if err == nil {
		fmt.Println(string(data))
	}

	// Obtener API key y modelo
	apiKey, err := config.GetAPIKey()
	if err != nil {
		ui.PrintError(err.Error())
		logger.Error("API key no configurada: %v", err)
		os.Exit(1)
	}

	model, err := config.GetModel()
	if err != nil {
		ui.PrintError("Error obteniendo modelo: " + err.Error())
		logger.Error("Error obteniendo modelo: %v", err)
		os.Exit(1)
	}

	ui.PrintInfo(fmt.Sprintf("\nModelo: %s", model))
	ui.PrintInfo(fmt.Sprintf("Conversación: %s\n", title))

	// Crear cliente y continuar conversación
	client := chat.NewClient(apiKey, model)
	handleInteractiveChatWithMessages(client, model, messages, title, timestamp)
}



