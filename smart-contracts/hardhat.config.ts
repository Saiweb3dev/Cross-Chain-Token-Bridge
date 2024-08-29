import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";
import "hardhat-deploy";
import type { NetworkUserConfig } from "hardhat/types";

const chainIds = {
  sepolia: 11155111,
  polygon_amoy:80002,
  hardhat: 31337,
  tBNB: 97,
}

const config: HardhatUserConfig = {
  defaultNetwork: "hardhat",
  networks: {
    hardhat: {
      chainId: chainIds.hardhat,
      allowUnlimitedContractSize: true
    },
    sepolia: {
      chainId: chainIds.sepolia,
      url: "https://sepolia.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
      accounts: process.env.PRIVATE_KEY !== undefined ? [process.env.PRIVATE_KEY] : [],
    },
    polygon_amoy: {
      chainId: chainIds.polygon_amoy,
      url: "https://polygon-mumbai.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161",
      accounts: process.env.PRIVATE_KEY !== undefined ? [process.env.PRIVATE_KEY] : [],
    },
  },
  solidity: {
    version:"0.8.19",
    settings:{
      optimizer:{
        enabled:true,
        runs:800
      }
    }
  },
  typechain: {
    outDir: "typechain",
    target: "ethers-v6",
  },
  namedAccounts: {
    deployer: {
      default: 0,
    },
  },
};

export default config;
