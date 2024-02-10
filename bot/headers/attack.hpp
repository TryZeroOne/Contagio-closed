#ifndef ATTACK
#define ATTACK

#include <map>
#include <string>

#include "include.hpp"

struct AttackStruct {
  string Name;
  string Ip;

  unsigned short int Duration;
  unsigned short int Port;
  unsigned short int ID;
};

class AttackReg {
 protected:
  using FunctionType = void (*)(types);

 public:
  void AddMethod(const string& name, FunctionType func);
  void LaunchMethod(const string& name, types);

 private:
  map<string, FunctionType> functions;
};

class Attack {
 private:
  bool* stop;

 public:
  AttackStruct as;

  Attack(AttackStruct as) {
    this->as = as;
  }
  Attack() {
  }

  void SetStop(bool* value) {
    this->stop = value;
  }
  bool* GetStop() {
    return this->stop;
  }

  void Start();
  void Stop();
};

extern int AttacksC;

#endif