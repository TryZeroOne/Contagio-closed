#ifndef UTILS
#define UTILS
#include <string>

unsigned int RandomInt(int min, int max);
void RandomBuf(char* buf, const int buffer_size);
std::string GetArch();

void DynamicName(char* argv[]);

#ifdef KILLER
bool IsInteger(const std::string str);
#endif

#endif