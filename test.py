#!/usr/bin/env python3
"""
Test script for all Propuestas API endpoints
"""

import requests
from _google_config import gcloud_token
from typing import LiteralString
from _update_static import update_static


# API base URL
BASE_URL = "http://localhost:8001"


def test_health_check():
    """Test health check endpoint"""
    print("🧪 Testing health check endpoint...")

    response = requests.get(f"{BASE_URL}/")

    if response.status_code == 200:
        data = response.json()
        print("✅ Health check successful!")
        print(f"   Status: {data['status']}")
        print(f"   Service: {data['service']}")
        print(f"   Version: {data['version']}")
        return True
    else:
        print(f"❌ Health check failed: {response.status_code}")
        print(f"   Error: {response.text}")
        return False


def test_gcr():
    """Test GCR endpoint"""
    print("🧪 Testing GCR endpoint...")
    url: LiteralString = "https://propuestas-api.or-gm.com/"

    payload: dict[str, str] = {
        "title": "Propuesta de Servicios",
        "subtitle": "Sistema de Transmisión Eléctrica",
        "prompt": """
        vamos a crear una propuesta nueva. GESTION de imdemnizacion por derecho de paso
         de linea de transmision de alta tesion. costo 4860 USD avacne 60% el resto con 
         la resolucion alcance levatamiento, acompanamiento, representacion, solicitud de servicio, 
         preparacion de expediente, ubicacion de torres, asesoria de volor de terreno o servidumbre., 
         seguimiennto, logistica, tiempo de servicio de 8 a 26 semanas dependiedo de la ETED
        """,
    }

    token = gcloud_token()
    resp_cloud: requests.Response = requests.post(
        url + "generate-text",
        json=payload,
        headers={"Authorization": f"Bearer {token}"},
    )
    return resp_cloud


def test_text_generation():
    """Test text generation endpoint"""
    print("🧪 Testing text generation endpoint...")

    payload: dict[str, str] = {
        "title": "Propuesta de iNSTALCION ELECTRICA",
        "subtitle": "MINISO GALERIA 360",
        "prompt": """
        vamos a crear una propuesta nueva. iNSTALACION ELECTRICA DE CANALIZACIONES EMT, PVC, CANALIZACIONE DE DATA, SUMINISTRO DE TRANSFORMADOR, EQUIPOS Y MATERIALES
        CANALIZACIONE DE BOCINA, INSTALCION DE TRANFORMADORT SECO, CABELADO Y CONEXION DE LUMINARIAS, INSTALCION DE LUMIANRIAS, PANELES ELECTRICOS
        SISTEMA DE CONTROL DE ILUMINACION, CAJA, LETROROS, LUCES DE EMERGENCIA. CABLEADO UPT DE DATA. TODO CON UN COSTO DE RD$. 600.000.00. TIEMPO DE EJECUCION DE 4 SEMANAS.
        AVANCE DE 70%.
        """,
        "model": "gpt-5-chat-latest",
    }

    response = requests.post(f"{BASE_URL}/generate-text", json=payload)

    if response.status_code == 200:
        data = response.json()
        print("✅ Text generation successful!")
        print(f"   Proposal ID: {data['id']}")
        print(f"   Created at: {data['created_at']}")
        print(f"   MD URL: {data['md_url']}")
        return data["id"]
    else:
        print(f"❌ Text generation failed: {response.status_code}")
        print(f"   Error: {response.text}")
        return None


def test_html_generation_from_id(proposal_id):
    """Test HTML generation from existing proposal ID"""
    print(f"\n🧪 Testing HTML generation from proposal ID: {proposal_id}...")

    payload = {"proposal_id": proposal_id, "model": "gpt-5-chat-latest"}

    response = requests.post(f"{BASE_URL}/generate-html", json=payload)

    if response.status_code == 200:
        data = response.json()
        print("✅ HTML generation from ID successful!")
        print(f"   Proposal ID: {data['id']}")
        print(f"   HTML URL: {data['html_url']}")
        print(f"   MD URL: {data['md_url']}")
        return True
    else:
        print(f"❌ HTML generation from ID failed: {response.status_code}")
        print(f"   Error: {response.text}")
        return False


def test_pdf_generation(proposal_id):
    """Test PDF generation from existing proposal ID"""
    print(f"\n🧪 Testing PDF generation from proposal ID: {proposal_id}...")

    payload = {"proposal_id": proposal_id}

    response = requests.post(f"{BASE_URL}/generate-pdf", json=payload)

    if response.status_code == 200:
        data = response.json()
        print("✅ PDF generation successful!")
        print(f"   Proposal ID: {data['id']}")
        print(f"   PDF URL: {data['pdf_url']}")
        return True
    else:
        print(f"❌ PDF generation failed: {response.status_code}")
        print(f"   Error: {response.text}")
        return False


def test_list_proposals():
    """Test listing all proposals"""
    print("\n🧪 Testing proposals list...")

    response = requests.get(f"{BASE_URL}/proposals")

    if response.status_code == 200:
        data = response.json()
        print("✅ Proposals list successful!")
        print(f"   Found {len(data)} proposals")
        for proposal in data:
            print(f"   - {proposal['id']}: {proposal['title']}")
            print(f"     HTML URL: {proposal.get('html_url', 'N/A')}")
            print(f"     PDF URL: {proposal.get('pdf_url', 'N/A')}")
            print(f"     Size HTML: {proposal.get('size_html', 0)} bytes")
            print(f"     Size PDF: {proposal.get('size_pdf', 0)} bytes")
        return True
    else:
        print(f"❌ Proposals list failed: {response.status_code}")
        print(f"   Error: {response.text}")
        return False


def test_search_proposals():
    """Test searching proposals by tokens in title/subtitle"""
    print("\n🧪 Testing proposals search...")
    q = "inter 69"
    response = requests.get(f"{BASE_URL}/proposals/search", params={"q": q})
    if response.status_code == 200:
        data = response.json()
        print("✅ Proposals search successful!")
        print(f"   Query: {q}")
        print(f"   Found {len(data)} proposals")
        for proposal in data:
            print(
                f"   - {proposal['id']}: {proposal['title']} | {proposal['subtitle']}"
            )
        return True
    else:
        print(f"❌ Proposals search failed: {response.status_code}")
        print(f"   Error: {response.text}")
        return False


def test_download_proposal(proposal_id, file_type):
    """Test downloading a proposal file (html or pdf)"""
    print(f"\n🧪 Testing download proposal {file_type} for ID: {proposal_id}...")

    response = requests.get(f"{BASE_URL}/proposal/{proposal_id}/{file_type}")

    with open(f"test/{proposal_id}.{file_type}", "wb") as f:
        f.write(response.content)

    if response.status_code == 200:
        print(f"✅ Download {file_type} successful!")
        print(f"   Content-Type: {response.headers.get('content-type', 'N/A')}")
        print(f"   Content-Length: {len(response.content)} bytes")
        print(
            f"   Content-Disposition: {response.headers.get('content-disposition', 'N/A')}"
        )
        return True
    else:
        print(f"❌ Download {file_type} failed: {response.status_code}")
        print(f"   Error: {response.text}")
        return False


def test_modify_proposal(proposal_id):
    """Test modifying an existing proposal"""
    print(f"\n🧪 Testing modify proposal for ID: {proposal_id}...")

    payload = {
        "title": "Propuesta de Servicios Modificada",
        "subtitle": "Sistema de Transmisión Eléctrica - Dark Mode",
        "prompt": "agrega una nota de que la desicion final la tiene EDESUR y en su caso elevado la superintendencia",
        "model": "gpt-5-chat-latest",
    }

    response = requests.put(f"{BASE_URL}/proposal/{proposal_id}", json=payload)

    if response.status_code == 200:
        data = response.json()
        print("✅ Modify proposal successful!")
        print(f"   Proposal ID: {data['id']}")
        print(f"   MD URL: {data['md_url']}")
        print(f"   Created at: {data['created_at']}")
        return True
    else:
        print(f"❌ Modify proposal failed: {response.status_code}")
        print(f"   Error: {response.text}")
        return False


def test_normal_workflow() -> None:
    # Test 1: Health check
    test_health_check()

    proposal_id: str | None = test_text_generation()

    if proposal_id:
        # Test 5: Generate HTML from existing proposal
        test_html_generation_from_id(proposal_id)

        # input("Press Enter to continue...")

        # Test 6: Generate PDF from existing proposal
        test_pdf_generation(proposal_id)

        # input("Press Enter to continue...")

        # # Test 7: Download MD file
        # test_download_proposal(proposal_id, "md")

        # #input("Press Enter to continue...")

        # # Test 8: Download HTML file
        # test_download_proposal(proposal_id, "html")

        # #input("Press Enter to continue...")

        # Test 9: Download PDF file
        test_download_proposal(proposal_id, "pdf")

        # input("Press Enter to continue...")

    # # Test 10: Modify existing proposal
    # test_modify_proposal(proposal_id)

    # # Test 11: List all proposals
    # test_list_proposals()

    # # Test 12: Search proposals
    # test_search_proposals()

    print("\n" + "=" * 60)
    print("✅ All tests completed!")


def main() -> None:
    # update_static()

    test_health_check()
    # proposal_id: str | None = test_text_generation()
    # test_html_generation_from_id(proposal_id)
    # test_download_proposal(proposal_id, "html")
    # update_static()
    # test_pdf_generation("0eed18b8")
    # test_download_proposal("0eed18b8", "pdf")
    # test_normal_workflow()


if __name__ == "__main__":
    main()

    # for filetype in ["html", "pdf", "md"]:
    #     test_download_proposal("fbf60f35", filetype)
