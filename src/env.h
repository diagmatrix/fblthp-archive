/**
 * Environment files management for the project
 * @author diagmatrix
 * @date 2024
 * @version 1.0
 */

#ifndef ENV_H
#define ENV_H
#include <string>
#include <cstdlib>
#include <iostream>
#include <filesystem>
#include <fstream>
#include <map>

namespace fs = std::filesystem;

/**
 * Tries to add variables from an environment file into the environment variables
 * @param env_filename Path for the environment file
 * @param defaults Default variables for the environment
 */
void load_env(const std::string& env_filename, const std::map<std::string, std::string>& defaults) {
    std::map<std::string, std::string> vars = defaults;
    
    if (fs::exists(env_filename) && fs::is_regular_file(env_filename)) {
        std::ifstream env_file(env_filename);
        std::string line;
        if (env_file.is_open()) {
           while(std::getline(env_file, line)) {
                // Ignore empty lines and comments
                if (line.empty() || line[0] == '#') {
                    continue;
                }
                
                size_t delimiterPos = line.find('=');
                
                // Ignore lines without an equal sign
                if (delimiterPos == std::string::npos) {
                    continue;
                }
                vars.insert_or_assign(line.substr(0, delimiterPos), line.substr(delimiterPos + 1));
           }
        }
    } else {
        std::cout << "Environment file could not be loaded. Using defaults..." << std::endl;
    }

    for (const auto& env_var: vars) {
        setenv(env_var.first.c_str(), env_var.second.c_str(), 1); // 1 = Override if defined
    }
}

#endif //ENV_H