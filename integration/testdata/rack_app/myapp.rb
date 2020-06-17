require 'bundler'
Bundler.require(:default)

class SimpleAdapter
  def call(env)
    [
      200,
      { 'Content-Type' => 'text/plain' },
      ["Hello world!"]
    ]
  end
end

Thin::Server.start('0.0.0.0', ENV["PORT"]) do
  use Rack::CommonLogger

  map '/' do
    run SimpleAdapter.new
  end
end
