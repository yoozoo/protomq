<?php
/**
 * @var Goridge\RelayInterface $relay
 */
use Spiral\Goridge;
use Spiral\RoadRunner;

ini_set('display_errors', 'stderr');
require 'vendor/autoload.php';
include 'GPBMetadata\Test.php';
include 'TypedQueue.php';

$rr = new RoadRunner\Worker(new Spiral\Goridge\StreamRelay(STDIN, STDOUT));

while ($body = $rr->receive($context)) {
    try {
    	$msg = new TypedQueue();
    	$msg->mergeFromString($body);
        // handle $msg
        

        $rr->send("", (string)$context);
    } catch (\Throwable $e) {
        $rr->error((string)$e);
    }
}
