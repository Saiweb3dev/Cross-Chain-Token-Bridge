import hre from "hardhat"

async function main(){
  const Token_Factory = await hre.ethers.getContractFactory("CCIP_Token");
  const Token = await Token_Factory.deploy("CCIP_Token","CCT",100);

  const contractAddress:string = typeof Token.target === 'string' ? Token.target : Token.target.toString();

  console.log(`The verifierETH Contract deployed at ${contractAddress}`)

  const abi = Token_Factory.interface.formatJson();
  const abiFormated = JSON.parse(abi);
  // console.log(abiFormated)
  // work in here of sending the abi properly
  await hre.deployments.save("CCIP_Token",{
    abi:abiFormated,
    address:contractAddress,
  })
}

main()
 .then(() => process.exit(0)) // Exit with success status code if deployment is successful
 .catch((error) => {
    console.error(error); // Log any errors that occur during deployment
    process.exit(1); // Exit with error status code if deployment fails
 });