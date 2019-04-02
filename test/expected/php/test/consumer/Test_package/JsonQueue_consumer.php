<?php
namespace Test_package;

use Spiral\Goridge;
use Spiral\RoadRunner;

ini_set('display_errors', 'stderr');
require 'vendor/autoload.php';

abstract class JsonQueue_consumer
{
    protected $rr;

    public function __construct()
    {
        $this->rr = new RoadRunner\Worker(new Spiral\Goridge\StreamRelay(STDIN, STDOUT));
    }

    public function run()
    {
        while ($body = $this->rr->receive($context)) {
            try {
                $this->handle_msg($body);

                $rr->send("", (string) $context);
            } catch (\Throwable $e) {
                $rr->error((string) $e);
            }
        }
    }

    abstract protected function handle_msg( $msg);
}
