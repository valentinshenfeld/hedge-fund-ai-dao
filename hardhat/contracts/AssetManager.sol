// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "../interfaces/IStrategyAdapter.sol";

contract AssetManager is AccessControl {
    using SafeERC20 for IERC20;

    bytes32 public constant EXECUTOR_ROLE = keccak256("EXECUTOR_ROLE"); // Role for CRE
    bytes32 public constant GOVERNOR_ROLE = keccak256("GOVERNOR_ROLE");

    event InvestmentExecuted(address indexed token, uint256 amount, string strategy);

    constructor(address _governor) {
        _grantRole(DEFAULT_ADMIN_ROLE, _governor);
        _grantRole(GOVERNOR_ROLE, _governor);
    }

    /**
     * @notice Главная функция инвестирования. 
     * Может быть вызвана ТОЛЬКО верифицированным Workflow DON (через Forwarder).
     */
    function invest(
        address token,
        uint256 amount,
        address adapter,
        bytes calldata strategyData
    ) external onlyRole(EXECUTOR_ROLE) {
        require(amount > 0, "Amount must be > 0");
        
        // 1. Проверки безопасности (например, Max Drawdown check)
        //...
        
        // 2. Аппрув токенов адаптеру стратегии (например, Aave или Uniswap)
        IERC20(token).forceApprove(adapter, amount);
        
        // 3. Выполнение входа в позицию
        IStrategyAdapter(adapter).enterPosition(token, amount, strategyData);
        
        emit InvestmentExecuted(token, amount, "StandardEntry");
    }

    // Функция для экстренного вывода средств (только Governor)
    function emergencyWithdraw(address token, address to) external onlyRole(GOVERNOR_ROLE) {
        uint256 balance = IERC20(token).balanceOf(address(this));
        IERC20(token).safeTransfer(to, balance);
    }
}