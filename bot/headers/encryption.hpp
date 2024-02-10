#ifndef ENCRYPTION
#define ENCRYPTION
#include <string>

#include "include.hpp"

void Sanitize(string& str);
extern string Decrypt(string __input);
#ifdef DEVMODE
void printBytes(const string& str);
#endif

#endif