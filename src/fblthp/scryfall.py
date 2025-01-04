from time import sleep
from typing import TypedDict, List

import requests
import sqlite3

# ----------------------------------------------------------------------------------------------------------------------
# SET
class MTGSet(TypedDict):
    """
    Represents the Scryfall data from a set
    """
    code: str
    name: str
    search_uri: str
    released_at: str
    card_count: int
    icon_svg_uri: str

def set_from_response(raw_set: dict) -> MTGSet:
    """
    Creates a MTGSet object from a Scryfall response set
    :param raw_set: Scryfall response set
    :return: MTGSet with the data from the Scryfall response
    """
    return MTGSet(
        code=raw_set.get("code", "UNKNOWN"),
        name=raw_set.get("name", "UNKNOWN"),
        search_uri=raw_set.get("search_uri", "UNKNOWN"),
        released_at=raw_set.get("released_at", "UNKNOWN"),
        card_count=raw_set.get("card_count", -1),
        icon_svg_uri=raw_set.get("icon_svg_uri", "UNKNOWN"),
    )

# ----------------------------------------------------------------------------------------------------------------------
# SCRYFALL CLIENT
class Scryfall:
    """
    Wrapper around the Scryfall API.
    """

    url = "https://api.scryfall.com/"
    insert_query = "INSERT INTO scryfall_history (url, headers, response_code, error) VALUES (?, ?, ?, ?)"

    def __init__(self, db: sqlite3.Connection, headers: dict = None):
        self.db = db
        if headers:
            self.headers = headers
        else:
            self.headers = {
                "Content-Type": "application/json",
                "User-Agent": "fblthp-archive-1.0"
            }

    def get_sets(self, digital: bool = False) -> List[MTGSet]:
        """
        Retrieves all scryfall sets.
        :param digital:
        :return: A list of MTGSet objects
        """

        sets = []
        sets_raw = self._request("sets")
        for sets_page in sets_raw:
            for set_raw in sets_page["data"]:
                if digital or (not set_raw.get("digital", True)):
                    sets.append(set_from_response(set_raw))

        return sets

    def _request(self, endpoint: str, full: bool = False) -> List[dict]:
        """
        Performs a request to the Scryfall API and logs it into the database.
        :param endpoint: Endpoint URL segment of the Scryfall API.
        :param full: If the endpoint segment is the full URL.
        :return: Response from the Scryfall API
        """

        url = f"{self.url}{endpoint}" if not full else endpoint

        sleep(0.1) # 100 ms wait before request
        response = requests.get(url, headers=self.headers)

        code = response.status_code
        error = response.json().get("error")
        self.db.execute(self.insert_query, (url, str(self.headers), code, error))
        self.db.commit()

        response_content = response.json()
        responses = [response_content]
        if response_content.get("has_more", False):
            responses.extend(self._request(response_content["next_page"], True))

        return responses
