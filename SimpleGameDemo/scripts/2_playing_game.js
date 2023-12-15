// We require the Hardhat Runtime Environment explicitly here. This is optional
// but useful for running the script in a standalone fashion through `node <script>`.
//
// You can also run a script with `npx hardhat run <script>`. If you do that, Hardhat
// will compile your contracts, add the Hardhat Runtime Environment's members to the
// global scope, and execute the script.
const hre = require("hardhat");
const { expect } = require("chai");
const configFile = process.cwd() + "/scripts/config.json";
const jsonfile = require('jsonfile');

const PreReady = 0;  // 0: 预准备状态
const Ready = 1;     // 1：准备状态
const Active = 2;    // 2: 游戏中
const Over = 3;       // 3:游戏结束

let TableId = 0; // 记录tableId
let game;

// 参与游戏
async function joinTable(TableId, player) {
  // 参与游戏
  await game.connect(player).joinTable(TableId);
}

// 查看玩家信息
async function getPlayerInfo(TableId, playerAddr) {
  console.log("getPlayerInfo=======");
  info = await game.getPlayerInfo(TableId, playerAddr);
  console.log("       player:",playerAddr, ", isInTable:", info.isInTable, ", cards:", info.cards);
}

// 查看游戏桌信息
async function getTableInfo(TableId) {
  console.log("getTableInfo, tableID:", TableId);
  table = await game.getTableInfo(TableId);
  console.log("       state:",table.state, ", round:", table.round, ", battleCards:", table.battleCards, ", playerAddrList:", table.playerAddrList);
  return table
}

async function main() {
  let config = await jsonfile.readFileSync(configFile);

  const [deployer, player1, player2] = await ethers.getSigners();

  // game合约
  const RockPaperScissorsGame = await ethers.getContractFactory("RockPaperScissorsGame");
  game = RockPaperScissorsGame.attach(config.ethSeries.RockPaperScissorsGame);
  
  console.log("Start to create game table");
  tx = await game.createTable();
  console.log("Create game table successful:", tx);

  console.log("Player start to join game table, tableID:", TableId);
  await joinTable(TableId, player1);
  await joinTable(TableId, player2);
  console.log("Player join game table successful");

  console.log("player2 readyGame");
  await game.connect(player2).readyGame(TableId);

  console.log("player1 startGame");
  await game.connect(player1).startGame(TableId);

  // 获取玩家将的卡牌
  await getPlayerInfo(TableId, player1.address);
  await getPlayerInfo(TableId, player2.address);


  tableInfo = await getTableInfo(TableId);
  // 游戏状态
  expect(tableInfo.state).to.equal(Active);
  expect(tableInfo.round).to.equal(0);  
  expect(tableInfo.playerAddrList.length).to.equal(2);  // 玩家人数
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main()
.catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
