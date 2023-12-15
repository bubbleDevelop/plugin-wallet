// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BehalfSignature {

    address private _issuerAddress;
    uint256 private _lineOfCredit;

    constructor() { _issuerAddress = msg.sender; }

    // 设置游戏运营地址，此地址上必须要有足够支付游戏过程中 gas 费的 token
   function setIssuer(address issuerAddress) public returns (bool success){
        require(msg.sender == _issuerAddress, "Only issuer can set new issuer");
        _issuerAddress = issuerAddress;
        return true;
    }
    
    // 获取游戏运营地址
    function issuer() external view returns(address){
        return _issuerAddress;
    }

    // 设置玩家的初始授信额度
    function setLineOfCredit(uint256 lineOfCreditAmount) public returns (bool success){
        require(msg.sender == _issuerAddress, "Only issuer can set line of credit");
        _lineOfCredit = lineOfCreditAmount;
        return true;
    }

    // 查询授信额度
    function lineOfCredit() external view returns (uint256) {
        return _lineOfCredit;
    }
}