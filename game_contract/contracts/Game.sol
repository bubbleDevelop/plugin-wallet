pragma solidity ^0.8.0;

contract Game {

    address private _issuerAddress;
    uint256 private _lineOfCredit;
    uint256 private _position;

   function setIssuer(address issuerAddress) public returns (bool success){
        _issuerAddress = issuerAddress;
        return true;
    }
    
    function issuer() external view returns(address){
        return _issuerAddress;
    }

    function setLineOfCredit(uint256 lineOfCreditAmount) public returns (bool success){
        _lineOfCredit = lineOfCreditAmount;
        return true;
    }

    function lineOfCredit() external view returns (uint256) {
        return _lineOfCredit;
    }

    function movePlayer(uint256 postionMove) public returns (bool success){
        _position = _position + postionMove; 
        return true;
    }

    function position() external view returns (uint256){
        return _position;
    }
}
