#ifndef METHODS
#define METHODS

#include <string>

#include "config.hpp"
#include "include.hpp"

void TcpmixMethod(types);
void UdpmixMethod(types);
void Launch(types, void (*method)(types));

#endif