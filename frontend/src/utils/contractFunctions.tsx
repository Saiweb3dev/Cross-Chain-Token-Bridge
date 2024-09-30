export const contractFunctions = {
  Token: ['mint', 'burn'],
  Vault: ['deposit', 'withdraw'],
  Router: ['swap', 'addLiquidity']
};

export type ContractType = keyof typeof contractFunctions;