import logging
import requests
import sqlite3
from time import sleep
from typing import TypedDict, Literal, get_args, Any, List
from dataclasses import dataclass

from pandas.core.arrays.arrow import ListAccessor


# ----------------------------------------------------------------------------------------------------------------------
# PAGED RESPONSE SCHEMA
class ResponseSchema(TypedDict):
    """
    Response schema for pages
    """
    total: int
    pages: list[dict]

# ----------------------------------------------------------------------------------------------------------------------
# SCRYFALL CLIENT
class ScryfallClient:
    """
    Wrapper around the Scryfall API.
    """

    url = "https://api.scryfall.com/"
    insert_query = "INSERT INTO scryfall_history (url, headers, response_code, error) VALUES (?, ?, ?, ?)"
    logger = logging.getLogger("Scryfall")

    def __init__(self, db: sqlite3.Connection, headers: dict = None):
        self.db = db
        if headers:
            self.headers = headers
        else:
            self.headers = {
                "Content-Type": "application/json",
                "User-Agent": "fblthp-archive-1.0"
            }

    def __send_request(self, url: str):
        """
        Sends a request to the Scryfall API.
        :param url: The Scryfall API URL
        :return: The response from the Scryfall API
        """
        sleep(0.1)  # 100 ms wait before request
        self.logger.info(f"Sending request to {url}")
        response = requests.get(url, headers=self.headers)

        code = response.status_code
        if code != 200:
            self.logger.warning(f"Request failed with code {code}")
        else:
            self.logger.info(f"Request successful")

        error = None if code == 200 else response.text
        self.db.execute(self.insert_query, (url, str(self.headers), code, error))
        self.db.commit()

        return response

    def request(self, url_segment: str) -> ResponseSchema:
        """
        Performs a request to the Scryfall API and logs it into the database.
        :param url_segment: Endpoint URL segment of the Scryfall API.
        :return: Response from the Scryfall API
        """

        self.logger.info(f"Requesting endpoint: {url_segment}")
        url = f"{self.url}{url_segment}"
        response = self.__send_request(url)
        responses: ResponseSchema = {
            "total": 1,
            "pages": [response.json()]
        }

        response_content = response.json()
        while response_content.get("has_more", False):
            response = self.__send_request(response_content["next_page"])
            responses["total"] += 1
            responses["pages"].append(response.json())
            response_content = response.json()

        self.logger.info(f"Retrieved {responses["total"]} pages")
        return responses

    def download(self, content_url: str, file_path: str):
        """
        Downloads a file from the Scryfall API.
        :param content_url: Content URL of the Scryfall API.
        :param file_path: Path to the file to save the download.
        """

        self.logger.info(f"Downloading file from {content_url}")
        response = self.__send_request(content_url)
        if response.ok:
            with open(file_path, "wb") as file:
                file.write(response.content)
            self.logger.info(f"Downloaded file {file_path}")
        else:
            self.logger.warning(f"Download failed with code {response.status_code}")

# ----------------------------------------------------------------------------------------------------------------------
# SET
SetType = Literal[
    "core", "expansion", "masters", "alchemy", "masterpiece", "arsenal", "from_the_vault", "spellbook", "premium_deck",
    "duel_deck", "draft_innovation", "treasure_chest", "commander", "planechase", "archenemy", "vanguard", "funny",
    "starter", "box", "promo", "token", "memorabilia", "minigame", "unknown"
]

def to_set_type(s: str) -> SetType:
    """
    Converts a string to a set type.
    :param s: String to convert
    :return: A set type converted from string (unknown if not a valid set type)
    """
    valid_types = get_args(SetType)

    return s if s in valid_types else "unknown"

@dataclass
class MTGSet:
    """
    Represents the Scryfall data from a set
    """
    code: str  # Set code
    name: str  # Set name
    set_type: SetType  # Set type
    digital: bool  # Whether set is digital only
    search_uri: str  # Set cards search URI
    released_at: str  # Set released date
    card_count: int  # Set unique printings count
    icon_svg_uri: str  # Set icon svg URI

    @classmethod
    def from_dict(cls, d: dict) -> "MTGSet":
        """
        Create a MTGSet from a dictionary (the Scryfall JSON response)
        :param d: Dictionary containing the MTGSet
        :return: A new MTGSet
        """
        return MTGSet(
            code=d.get("code", "UNKNOWN"),
            name=d.get("name", "UNKNOWN"),
            set_type=to_set_type(d.get("set_type", "")),
            digital=d.get("digital", False),
            search_uri=d.get("search_uri", "UNKNOWN"),
            released_at=d.get("released_at", "UNKNOWN"),
            card_count=d.get("card_count", -1),
            icon_svg_uri=d.get("icon_svg_uri", "UNKNOWN"),
        )

    def to_list(self) -> List[Any]:
        """
        :return: The MTGSet as a list (name, code, set_type, digital, released_at, card_count, search_uri, icon_uri)
        """
        return [
            self.name,
            self.code,
            self.set_type,
            self.digital,
            self.released_at,
            self.card_count,
            self.search_uri,
            self.icon_svg_uri
        ]
