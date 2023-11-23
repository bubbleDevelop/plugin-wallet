const configFile = process.cwd() + "/config.json";
const jsonfile = require('jsonfile')
const HDWalletProvider = require('@truffle/hdwallet-provider');
var Game = artifacts.require("Game");

module.exports = async function(deployer) {
    let config = await jsonfile.readFile(configFile);
  
    // vote contract
    var provider = new HDWalletProvider(config.bubbledev.mnemonic, 'ws://192.168.31.115:18002');

    let gameAddress = await deployer.deploy(Game);
    console.log('game contract address:', gameAddress.address);
    config.game.address = gameAddress.address;
  
    console.log("deploy successful");
  
    await jsonfile.writeFile(configFile, config, {spaces: 2});
    console.log("deploy Game contract config success");
  };