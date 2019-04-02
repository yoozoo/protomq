<?php
namespace Test_package;

use Spiral\Goridge;

class JsonQueueProducer
{
    /**
     * @var Goridge\RPC
     */
    private $rpc;

    public function __construct($addr = "127.0.0.1", $port = 8080)
    {
        $this->rpc = new Goridge\RPC(new Goridge\SocketRelay($addr, $port));
    }

    public function send($data)
    {
        $payload = array(
            "topic" => "test",
            "content" => $data,
        );
        $this->rpc->call("Sender.Send", $payload);
    }
}
