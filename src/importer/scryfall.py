import logging
import os
from typing import List, TYPE_CHECKING

from models import ScryfallClient, MTGSet, SetType

if TYPE_CHECKING:
    import sqlite3

# ----------------------------------------------------------------------------------------------------------------------
# LOGGER
logger = logging.getLogger("Parser")

# ----------------------------------------------------------------------------------------------------------------------
# CONSTANTS
INSERT_SET_QUERY = "INSERT INTO mtg_set (name, code, set_type, digital, released_at, card_count, search_uri, icon_uri) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"

# ----------------------------------------------------------------------------------------------------------------------
# SET FUNCTIONS
def get_sets(client: ScryfallClient) -> List[MTGSet]:
    """
    Returns a list of all Magic: The Gathering sets
    :param client: Scryfall client
    :return: List of Magic: The Gathering sets
    """
    logger.info("Retrieving sets")
    response = client.request("sets")["pages"]
    sets: List[MTGSet] = []
    logger.info("Parsing sets")
    for page in response:
        for set_raw in page["data"]:
            sets.append(MTGSet.from_dict(set_raw))

    return sets

def remove_sets(sets: List[MTGSet], set_types: List[SetType], digital: bool) -> List[MTGSet]:
    """
    Removes sets from a list of sets
    :param sets: List of sets
    :param set_types: Set types to remove
    :param digital: Whether to remove digital sets
    :return: Filtered list of sets
    """
    logger.info("Removing unwanted sets")
    new_sets: List[MTGSet] = []
    for set in sets:
        if set.set_type in set_types:
            logger.info(f"Removing set {set.name} [{set.code}] ({set.set_type})")
        elif digital and set.digital:
            logger.info(f"Removing set {set.name} [{set.code}] (digital)")
        else:
            new_sets.append(set)

    return new_sets

def download_icons(client: ScryfallClient, sets: List[MTGSet], folder: str) -> None:
    """
    Downloads icons from a list of sets
    :param client: Scryfall client
    :param sets: List of sets
    :param folder: Folder to download icons to
    """
    logger.info("Downloading icons")
    if not os.path.exists(folder):
        os.makedirs(folder)

    for set in sets:
        logger.info(f"Downloading icons for set {set.name} [{set.code}]")
        set_folder = os.path.join(folder, set.code)
        if not os.path.exists(set_folder):
            os.makedirs(set_folder)
        icon_name = os.path.join(set_folder, "icon.svg")
        if not os.path.exists(icon_name):
            client.download(set.icon_svg_uri, icon_name)
        else:
            logger.info("Icon already exists, skipping download")

def process_sets(
        db: "sqlite3.Connection",
        client: ScryfallClient,
        exclude_types: List[SetType],
        exclude_digital: bool,
        folder: str
) -> List[MTGSet]:
    """
    Processes scryfall sets, downloading the icons and inserting them in the database.
    :param db: Database connection to insert the scryfall sets into
    :param client: Scryfall client
    :param exclude_types: Set types to exclude
    :param exclude_digital: Whether to exclude digital sets
    :param folder: Folder to download icons to
    :return: The list of processed sets
    """
    logger.info("Processing scryfall sets")

    sets = get_sets(client)
    sets = remove_sets(sets, exclude_types, exclude_digital)
    download_icons(client, sets, folder)

    logger.info("Inserting scryfall sets into the database")
    set_list = [mtg_set.to_list() for mtg_set in sets]
    db.executemany(INSERT_SET_QUERY, set_list)
    db.commit()

    return sets
