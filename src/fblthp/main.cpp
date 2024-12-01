/**
 * Fblthp collection manager
 * @author diagmatrix
 * @date 2024
 * @version 1.0
 */
#include <filesystem>
#include <iostream>
#include <curl/curl.h>
#include <nlohmann/json.hpp>

namespace fs = std::filesystem;
using json = nlohmann::json;

const char* ALL_CARDS_URL = "https://api.scryfall.com/bulk-data/all-cards";
const char* ALL_CARDS_FILE = "all_cards.json";

size_t write_to_file(void* contents, size_t size, size_t nmemb, FILE* response_file) {
    size_t total_size = size * nmemb;
    fwrite(contents, size, nmemb, response_file);
    return total_size;
}

size_t write_to_string(void* contents, size_t size, size_t nmemb, std::string* response_string) {
    size_t total_size = size * nmemb;
    response_string->append((char*)contents, total_size);
    return total_size;
}

void set_curl_headers(CURL* curl, curl_slist* headers) {
    headers = curl_slist_append(headers, "Content-Type: application/json");
    headers = curl_slist_append(headers, "User-Agent: fblthp/1.0");
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, 10L);
}

int main() {
    // Initialize curl
    CURL* curl = curl_easy_init();
    if (!curl) {
        std::cout << "Error: Unable to initialize curl" << std::endl;
        return EXIT_FAILURE;
    }

    // Initialize string to write to
    std::string response_string;

    // Set curl options
    curl_easy_setopt(curl, CURLOPT_URL, ALL_CARDS_URL);
    curl_slist* headers = nullptr;
    set_curl_headers(curl, headers);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_to_string);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response_string);

    // Get all cards bulk data
    CURLcode res = curl_easy_perform(curl);
    if (res != CURLE_OK) {
        std::cout << "Error: " << curl_easy_strerror(res) << std::endl;
        return EXIT_FAILURE;
    }
    json response_json = json::parse(response_string);
    std::cout << response_json.dump(4) << std::endl;


    // Create file to write to
    // if (fs::exists(ALL_CARDS_FILE)) {  // TODO: Add option to force the download
    //     std::cout << "All cards have been downloaded" << std::endl;
    //     return EXIT_SUCCESS;
    // }
    //
    // FILE* file = fopen(ALL_CARDS_FILE, "w");
    // if (!file) {
    //     std::cout << "Error: Unable to create file" << std::endl;
    //     return EXIT_FAILURE;
    // }
    //
    // if (res != CURLE_OK) {
    //     std::cout << "Error: " << curl_easy_strerror(res) << std::endl;
    //     fclose(file);
    //     return EXIT_FAILURE;
    // }
    //
    // std::cout << "All cards have been downloaded" << std::endl;
    // fclose(file);
    curl_easy_cleanup(curl);
    curl_slist_free_all(headers);

    return EXIT_SUCCESS;
}