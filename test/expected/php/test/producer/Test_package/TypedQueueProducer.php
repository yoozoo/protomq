<?php
namespace Test_package;

use Spiral\Goridge;

class TypedQueueProducer
{
    /**
     * @var Goridge\RPC
     */
    private $rpc;

    public function __construct($addr = "127.0.0.1", $port = 8080)
    {
        $this->rpc = new Goridge\RPC(new Goridge\SocketRelay($addr, $port));
    }

    public function send(TypedQueue $data)
    {
        $payload = array(
            "topic" => "test_typed",
            "content" => $data->serializeToString(),
        );
        $this->rpc->call("Sender.Send", $payload);
    }
}
