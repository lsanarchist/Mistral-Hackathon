#!/usr/bin/env ruby

require 'webrick'
require 'json'

# Simple Ruby web server for demo purposes
class RubyDemoServer < WEBrick::HTTPServlet::AbstractServlet
  def do_GET(request, response)
    case request.path
    when '/'
      response.body = "Ruby Demo Server - TriageProf"
    when '/cpu-intensive'
      cpu_intensive(response)
    when '/memory-heavy'
      memory_heavy(response)
    when '/object-creation'
      object_creation(response)
    when '/json-processing'
      json_processing(response)
    when '/database-queries'
      database_queries(response)
    else
      response.status = 404
      response.body = "Not Found"
    end
  end

  private

  def cpu_intensive(response)
    start_time = Time.now
    # Simulate CPU work
    result = 0
    10000000.times do |i|
      result += i * i
    end
    response.body = "CPU intensive completed in #{Time.now - start_time} seconds. Result: #{result}"
  end

  def memory_heavy(response)
    # Simulate memory allocation
    arrays = []
    10000.times do
      arrays << Array.new(1000) { |i| i.to_s }
    end
    response.body = "Memory heavy completed. Allocated #{arrays.size} arrays"
  end

  def object_creation(response)
    # Simulate object creation
    objects = []
    5000.times do
      objects << { id: rand(1000), name: "Object_#{rand(1000)}", data: Array.new(100) { rand } }
    end
    response.body = "Object creation completed. Created #{objects.size} objects"
  end

  def json_processing(response)
    # Simulate JSON processing
    data = { users: [], products: [] }
    1000.times do
      data[:users] << { id: rand(1000), name: "User_#{rand(1000)}", email: "user_#{rand(1000)}@example.com" }
      data[:products] << { id: rand(1000), name: "Product_#{rand(1000)}", price: rand(100.0) }
    end
    response.body = "JSON processing completed. Generated #{data[:users].size} users and #{data[:products].size} products"
  end

  def database_queries(response)
    # Simulate database queries
    results = []
    1000.times do
      results << {
        id: rand(1000),
        name: "Record_#{rand(1000)}",
        created_at: Time.now.iso8601,
        updated_at: Time.now.iso8601
      }
    end
    response.body = "Database queries completed. Retrieved #{results.size} records"
  end
end

# Start the server
server = WEBrick::HTTPServer.new(Port: 4567)
server.mount "/", RubyDemoServer

puts "Ruby Demo Server running on http://localhost:4567"
puts "Endpoints:"
puts "- /cpu-intensive - CPU intensive operations"
puts "- /memory-heavy - Memory allocation heavy operations"
puts "- /object-creation - Object creation operations"
puts "- /json-processing - JSON processing operations"
puts "- /database-queries - Database query simulation"

# Run in a separate thread so we can handle Ctrl+C
Thread.new { server.start }

# Wait for interrupt
trap("INT") { server.shutdown }
sleep