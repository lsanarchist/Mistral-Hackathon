#!/usr/bin/env ruby

require 'json'
require 'net/http'
require 'uri'
require 'tempfile'
require 'fileutils'

class RubyStackprofPlugin
  def initialize
    @info = {
      "name" => "ruby-stackprof",
      "version" => "0.1.0",
      "sdkVersion" => "1.0",
      "capabilities" => {
        "targets" => ["url"],
        "profiles" => ["cpu", "memory", "object_allocation"]
      }
    }
  end

  def handle_rpc(request)
    method = request["method"]
    params = request["params"] || {}
    rpc_id = request["id"] || 0

    result = nil
    error_obj = nil

    begin
      case method
      when "rpc.info"
        result = @info
      when "rpc.validateTarget"
        result = validate_target(params)
      when "rpc.collect"
        result = collect(params)
      else
        error_obj = {
          "code" => -32601,
          "message" => "Method not found: #{method}"
        }
      end
    rescue => e
      error_obj = {
        "code" => -32603,
        "message" => e.message,
        "data" => e.backtrace&.first(5)&.join("\n")
      }
    end

    response = {
      "jsonrpc" => "2.0",
      "id" => rpc_id
    }

    if error_obj
      response["error"] = error_obj
    else
      response["result"] = result
    end

    response
  end

  def validate_target(target)
    return false unless target["type"] == "url"
    return false unless target["baseUrl"]&.start_with?("http://", "https://")
    true
  end

  def collect(request)
    target_url = request["target"]["baseUrl"]
    duration_sec = request["durationSec"] || 30
    out_dir = request["outDir"] || "./out"
    profiles = request["profiles"] || ["cpu", "memory", "object_allocation"]

    # Create output directory
    FileUtils.mkdir_p(out_dir)

    artifacts = []
    metadata = {
      "timestamp" => Time.now.utc.strftime('%Y-%m-%dT%H:%M:%SZ'),
      "durationSec" => duration_sec,
      "service" => "ruby-application",
      "scenario" => "production"
    }

    # Simulate collecting profiles from Ruby stackprof
    # In a real implementation, this would connect to a Ruby app with stackprof
    profiles.each do |profile_type|
      case profile_type
      when "cpu"
        file_path = File.join(out_dir, "cpu.stackprof")
        File.write(file_path, generate_sample_cpu_profile)
        artifacts << create_artifact("stackprof", "cpu", file_path)
      when "memory"
        file_path = File.join(out_dir, "memory.stackprof")
        File.write(file_path, generate_sample_memory_profile)
        artifacts << create_artifact("stackprof", "memory", file_path)
      when "object_allocation"
        file_path = File.join(out_dir, "object_allocation.stackprof")
        File.write(file_path, generate_sample_object_allocation_profile)
        artifacts << create_artifact("stackprof", "object_allocation", file_path)
      end
    end

    {
      "metadata" => metadata,
      "target" => request["target"],
      "artifacts" => artifacts
    }
  end

  private

  def create_artifact(kind, profile_type, path)
    {
      "kind" => kind,
      "profileType" => profile_type,
      "path" => path,
      "contentType" => "application/octet-stream"
    }
  end

  def generate_sample_cpu_profile
    <<~PROFILE
      ==1973==
      Mode: cpu(1000)
      Samples: 1000
      Thread: 1973
      ==1973==
      400  main
        main@/app/main.rb:5
      300  process_request
        process_request@/app/request_handler.rb:12
      200  parse_json
        parse_json@/app/json_parser.rb:8
      100  database_query
        database_query@/app/db.rb:15
    PROFILE
  end

  def generate_sample_memory_profile
    <<~PROFILE
      ==1974==
      Mode: object(1000)
      Samples: 1000
      Thread: 1974
      ==1974==
      500  String.new
        String.new@/app/main.rb:10
      300  Array.new
        Array.new@/app/request_handler.rb:20
      200  Hash.new
        Hash.new@/app/json_parser.rb:15
    PROFILE
  end

  def generate_sample_object_allocation_profile
    <<~PROFILE
      ==1975==
      Mode: object(1000)
      Samples: 1000
      Thread: 1975
      ==1975==
      600  User.new
        User.new@/app/models/user.rb:5
      250  Product.new
        Product.new@/app/models/product.rb:8
      150  Order.new
        Order.new@/app/models/order.rb:12
    PROFILE
  end
end

# Main execution
if __FILE__ == $0
  plugin = RubyStackprofPlugin.new

  # Read from stdin, write to stdout
  $stdin.each_line do |line|
    begin
      request = JSON.parse(line)
      response = plugin.handle_rpc(request)
      puts JSON.generate(response)
      $stdout.flush
    rescue JSON::ParserError => e
      $stderr.puts "Error parsing JSON: #{e.message}"
    rescue => e
      $stderr.puts "Error: #{e.message}"
    end
  end
end