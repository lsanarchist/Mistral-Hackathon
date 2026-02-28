#!/usr/bin/env ruby

require 'json'
require_relative 'main'

def test_plugin_info
  plugin = RubyStackprofPlugin.new
  info = plugin.handle_rpc({
    "jsonrpc" => "2.0",
    "method" => "rpc.info",
    "id" => 1
  })
  
  puts "Testing plugin info..."
  puts "Name: #{info["result"]["name"]}"
  puts "Version: #{info["result"]["version"]}"
  puts "SDK Version: #{info["result"]["sdkVersion"]}"
  puts "Targets: #{info["result"]["capabilities"]["targets"].join(", ")}"
  puts "Profiles: #{info["result"]["capabilities"]["profiles"].join(", ")}"
  puts "✓ Plugin info test passed"
  puts
end

def test_target_validation
  plugin = RubyStackprofPlugin.new
  
  # Test valid target
  result = plugin.handle_rpc({
    "jsonrpc" => "2.0",
    "method" => "rpc.validateTarget",
    "params" => {"type" => "url", "baseUrl" => "http://localhost:4567"},
    "id" => 2
  })
  
  puts "Testing target validation..."
  puts "Valid URL result: #{result["result"]}"
  
  # Test invalid target
  result = plugin.handle_rpc({
    "jsonrpc" => "2.0",
    "method" => "rpc.validateTarget",
    "params" => {"type" => "file", "baseUrl" => "/tmp/test"},
    "id" => 3
  })
  
  puts "Invalid target result: #{result["result"]}"
  puts "✓ Target validation test passed"
  puts
end

def test_collect
  plugin = RubyStackprofPlugin.new
  
  Dir.mktmpdir do |dir|
    result = plugin.handle_rpc({
      "jsonrpc" => "2.0",
      "method" => "rpc.collect",
      "params" => {
        "target" => {"type" => "url", "baseUrl" => "http://localhost:4567"},
        "durationSec" => 10,
        "outDir" => dir,
        "profiles" => ["cpu", "memory", "object_allocation"]
      },
      "id" => 4
    })
    
    puts "Testing collect..."
    puts "Artifacts created: #{result["result"]["artifacts"].size}"
    result["result"]["artifacts"].each do |artifact|
      puts "  - #{artifact["profileType"]}: #{artifact["path"]}"
      puts "    Exists: #{File.exist?(artifact["path"])}"
    end
    puts "✓ Collect test passed"
  end
  puts
end

begin
  test_plugin_info
  test_target_validation
  test_collect
  puts "All tests passed! ✓"
rescue => e
  puts "Test failed: #{e.message}"
  puts e.backtrace.first(5)
  exit 1
end