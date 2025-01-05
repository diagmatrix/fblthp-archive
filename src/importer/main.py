import sqlite3
import logging

from scryfall import process_sets
from models import ScryfallClient

logger = logging.getLogger("Card Extractor")

if __name__ == "__main__":
    logging.basicConfig(
        format="[%(asctime)s] %(levelname)s %(name)s: %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
        level=logging.INFO,
    )

    logger.info("Starting...")

    connection = sqlite3.connect("../../archive.db")
    client = ScryfallClient(connection)

    sets = process_sets(connection, client, ["funny", "memorabilia"], True, "../../data/")
    logger.info(f"Processed {len(sets)} sets")

    logger.info("Finished")