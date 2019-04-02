<?php
namespace Test_package;

use Exception;
use Spiral\Goridge;
use Spiral\RoadRunner;

class JsonQueueConsumer
{
    private $rr;
    private $handler;

    public function __construct()
    {
        $this->rr = new RoadRunner\Worker(new Goridge\StreamRelay(STDIN, STDOUT));
        $this->handler = function ($data) {
            throw new Exception("No handler registered.");
        };
    }

    public function register_handler($handler)
    {
        if (!is_callable($handler)) {
            throw new Exception("The 'handler' is not callable.");
        }
        $this->handler = $handler;
    }

    public function run()
    {
        while ($body = $this->rr->receive($context)) {
            try {
                call_user_func($this->handler, $body);

                $this->rr->send("", (string) $context);
            } catch (\Throwable $e) {
                $this->rr->error((string) $e);
            }
        }
    }
}
