import { deployments } from "hardhat";
const fs = require("fs");

const BACK_END_TOKEN_ABI_FILE = "../backend/contractDetails/tokenContractABI.json";
const BACK_END_VAULT_ABI_FILE = "../backend/contractDetails/vaultContractABI.json";
const BACK_END_MESSANGER_ABI_FILE = "../backend/contractDetails/messangerContractABI.json";

const main = async () => {
  if(process.env.UPDATE_BACK_END){
    console.log("Updating contract ABIs...")
    
    // Update CCIP_Token ABI
    const tokenContract = await deployments.get("CCIP_Token");
    await updateAbi(tokenContract.abi, BACK_END_TOKEN_ABI_FILE, "CCIP_Token");
    
    // Update CCIP_TokenVault ABI
    const vaultContract = await deployments.get("CCIP_TokenVault");
    await updateAbi(vaultContract.abi, BACK_END_VAULT_ABI_FILE, "CCIP_TokenVault");
    
    // Update CrossChain_Messanger ABI
    const messangerContract = await deployments.get("CrossChain_Messanger");
    await updateAbi(messangerContract.abi, BACK_END_MESSANGER_ABI_FILE, "CrossChain_Messanger");
  }
  else{
    console.log("NO Permission Given")
    console.log(process.env.UPDATE_BACK_END)
  }
}

async function updateAbi(contractABI: any, filePath: string, contractName: string){
  console.log(`Updating the ABI for ${contractName}`)
  const abiJson = JSON.stringify(contractABI, null, 2);
  fs.writeFileSync(filePath, abiJson);
  console.log(`${contractName} ABI updated successfully`)
}

main().then(() => process.exit(0)).catch((error) => {
  console.error(error)
  process.exit(1);
})