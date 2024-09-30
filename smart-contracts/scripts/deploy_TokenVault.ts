import hre from "hardhat"

async function main() {
  const tokenAddress = "0x03A07c5991B70497813F8Cb4c886F19e1A231d5c";

  const TokenVault_Factory = await hre.ethers.getContractFactory("contracts/CCIP_TokenVault.sol:CCIP_TokenVault");
  const TokenVault = await TokenVault_Factory.deploy(tokenAddress);

  const contractAddress: string = typeof TokenVault.target === 'string' ? TokenVault.target : TokenVault.target.toString();

  console.log(`The CCIP_TokenVault Contract deployed at ${contractAddress}`);

  const abi = TokenVault_Factory.interface.formatJson();
  const abiFormatted = JSON.parse(abi);

  await hre.deployments.save("CCIP_TokenVault", {
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