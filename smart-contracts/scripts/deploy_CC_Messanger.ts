import hre from "hardhat"

async function main() {
  // Addresses for Sepolia testnet
  const routerAddress = "0xD0daae2231E9CB96b94C8512223533293C3693Bf"; // Sepolia CCIP router
  const linkAddress = "0x779877A7B0D9E8603169DdbD7836e478b4624789";  // Sepolia LINK token
  const vaultAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3"; // Provided vault address

  const CrossChainMessanger_Factory = await hre.ethers.getContractFactory("contracts/CrossChain_Messanger.sol:CrossChain_Messanger");
  const CrossChainMessanger = await CrossChainMessanger_Factory.deploy(routerAddress, linkAddress, vaultAddress);

  const contractAddress: string = typeof CrossChainMessanger.target === 'string' ? CrossChainMessanger.target : CrossChainMessanger.target.toString();

  console.log(`The CrossChain_Messanger Contract deployed at ${contractAddress}`);

  const abi = CrossChainMessanger_Factory.interface.formatJson();
  const abiFormatted = JSON.parse(abi);

  await hre.deployments.save("CrossChain_Messanger", {
    abi: abiFormatted,
    address: contractAddress,
  });
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });