<?php
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: test.proto

namespace Test_package;

use Google\Protobuf\Internal\GPBType;
use Google\Protobuf\Internal\RepeatedField;
use Google\Protobuf\Internal\GPBUtil;

/**
 * Generated from protobuf message <code>test_package.TypedQueue</code>
 */
class TypedQueue extends \Google\Protobuf\Internal\Message
{
    /**
     * Generated from protobuf field <code>.test_package.Log data = 1;</code>
     */
    private $data = null;

    /**
     * Constructor.
     *
     * @param array $data {
     *     Optional. Data for populating the Message object.
     *
     *     @type \Test_package\Log $data
     * }
     */
    public function __construct($data = NULL) {
        \GPBMetadata\Test::initOnce();
        parent::__construct($data);
    }

    /**
     * Generated from protobuf field <code>.test_package.Log data = 1;</code>
     * @return \Test_package\Log
     */
    public function getData()
    {
        return $this->data;
    }

    /**
     * Generated from protobuf field <code>.test_package.Log data = 1;</code>
     * @param \Test_package\Log $var
     * @return $this
     */
    public function setData($var)
    {
        GPBUtil::checkMessage($var, \Test_package\Log::class);
        $this->data = $var;

        return $this;
    }

}

