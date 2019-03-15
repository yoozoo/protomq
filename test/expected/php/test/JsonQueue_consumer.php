<?php
/**
 * @var Goridge\RelayInterface $relay
 */
use Spiral\Goridge;
use Spiral\RoadRunner;

ini_set('display_errors', 'stderr');
require 'vendor/autoload.php';

$rr = new RoadRunner\Worker(new Spiral\Goridge\StreamRelay(STDIN, STDOUT));

while ($body = $rr->receive($context)) {
    try {
        // handle $body string

        $rr->send("", (string)$context);
    } catch (\Throwable $e) {
        $rr->error((string)$e);
    }
}
