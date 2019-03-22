<?php
/**
 * @var Goridge\RelayInterface $relay
 */
namespace Test_package;

use Spiral\Goridge;
use Spiral\RoadRunner;

ini_set('display_errors', 'stderr');
require 'vendor/autoload.php';
include 'GPBMetadata\Test.php';
include 'TypedQueue.php';

abstract class TypedQueue_consumer
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
                $msg = new TypedQueue();
                $msg->mergeFromString($body);
                $this->handle_msg($msg);

                $rr->send("", (string) $context);
            } catch (\Throwable $e) {
                $rr->error((string) $e);
            }
        }
    }

    abstract protected function handle_msg(TypedQueue $msg);
}
