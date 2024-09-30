export const contractFunctions = {
  Token: ['mint', 'burn', 'transfer', 'balanceOf', 'totalSupply', 'approve', 'transferFrom', 'allowance'],
  Vault: ['lockTokenInVault', 'releaseTokenInVault', 'owner'],
  Router: ['sendMessagePayLINK', 'getAllClientData', 'getLastReceivedMessageDetails', 'withdraw', 'withdrawToken'],
  Messenger: ['sendMessagePayLINK', 'getAllClientData', 'getLastReceivedMessageDetails', 'withdraw', 'withdrawToken']
};

export type ContractType = keyof typeof contractFunctions;