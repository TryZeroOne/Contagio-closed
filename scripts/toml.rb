def parse_toml(file_path)
    parsed_data = {}
  
    unless File.exist?(file_path)
      puts "Invalid path: #{file_path}"
      return parsed_data
    end
  
    File.open(file_path, 'r') do |file|
      table = parsed_data
  
      file.each_line do |line|
        line.strip!
  
        next if line.empty? || line.start_with?('#')
  
        if line.start_with?('[') && line.end_with?(']')
          table = parsed_data
          keys = line[1...-1].split('.')
          keys.each do |key|
            table[key] ||= {}
            table = table[key]
          end
        elsif line.include?('=')
          key, value = line.split('=', 2).map(&:strip)
          key = key.gsub(/"|'/, '')
          value = value.gsub(/"|'/, '')
  
          if value.match?(/^\d+$/)
            value = value.to_i
          elsif value.match?(/^\d+\.\d+$/)
            value = value.to_f
          elsif value.downcase == "true"
            value = true
          elsif value.downcase == "false"
            value = false
          end
  
          table[key] = sanitize(value)
        end
      end
    end
  
    parsed_data
  end
  
def sanitize(input)
  if input.is_a?(String)
    len = input.length
    result = ''
    i = 0
    while i < len
      if input[i] == '#'
        break
      else
        result += input[i]
        i += 1
      end
    end

    if result[-1] == ' '
      return result.chop
    end  

    return result
  elsif input == true
    return "true"
  elsif input == false
    return "false"
  else
    return "Unsupported input type"
  end
end