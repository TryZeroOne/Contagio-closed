#include "headers/encryption.hpp"

string Decrypt(string __input) {
  const string key = __input.substr(0, 128) + __input.substr(__input.length() - 128);
  string command = __input.substr(128, __input.length() - 256);
  string _result(command.length(), '\0');

  int i = 0;
  for (char c : command) {
    _result += (c - key[i % key.length()]);
    i++;
  }

  Sanitize(_result);
  return _result;
}

void Sanitize(string& str) {
  const unsigned char* bytes = reinterpret_cast<const unsigned char*>(str.c_str());
  unsigned long length = str.length();

  string new_string;

  for (unsigned long i = 0; i < length; i++) {
    if (static_cast<int>(bytes[i]) == 0) {
      continue;
    }

    new_string += static_cast<int>(bytes[i]);
  }

  str = new_string;
}