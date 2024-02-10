def hash(str)
    command = "go run contagio/enc/main.go hash " + str
    result = `#{command}`
    result
  end

def main
    if ARGV.length >= 2
      login = hash(ARGV[0])
      password = hash(ARGV[1])
  
      system("rm sqlite/database.db > /dev/null 2>&1")
      system("mkdir sqlite > /dev/null 2>&1")
      system("touch sqlite/database.db > /dev/null 2>&1")
      system("sqlite3 sqlite/database.db 'CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, login TEXT, password TEXT)'")
      system("sqlite3 sqlite/database.db 'CREATE TABLE IF NOT EXISTS allowed (id INTEGER PRIMARY KEY, ip TEXT)'")
      system("sqlite3 sqlite/database.db \"INSERT INTO users(login, password) VALUES('#{login}', '#{password}')\"")
  
      puts "\n=================================\n\t      Login\n=================================\n"
      puts login
  
      puts "\n=================================\n\t     Password\n=================================\n"
      puts password
    else
      puts "Invalid args!\n\nruby scripts/init.rb <LOGIN> <PASSWORD>\nExample:\nruby scripts/init.rb usr pass"
    end
end
  
main
  