// SPDX-License-Identifier: MIT
pragma solidity ^0.8.27;

/// @notice Minimal ATP interface exposing balance-related getters
interface IATPCore {
    function getClaimable() external view returns (uint256);
    function getAllocation() external view returns (uint256);
    function getClaimed() external view returns (uint256);
}

/// @title ATP Balance Helper
/// @notice Computes VESTED, UNVESTED, and CLAIMED balances for Aztec Token Positions (LATP / MATP)
contract ATPBalanceHelper {
    struct Balances {
        uint256 allocation; // allocation
        uint256 vested;     // claimable balance
        uint256 unvested;   // allocation - claimed - claimable
        uint256 claimed;    // already claimed amount
    }

    /// @notice Compute all balance types for a given ATP
    /// @param atp Address of the ATP (LATP or MATP) contract
    /// @return balances Struct containing vested, unvested, and claimed amounts
    function getBalances(address atp) external view returns (Balances memory balances) {
        IATPCore core = IATPCore(atp);

        uint256 claimable = core.getClaimable();
        uint256 allocation = core.getAllocation();
        uint256 claimed = core.getClaimed();

        balances.allocation = allocation;
        balances.vested = claimable;
        balances.claimed = claimed;
        balances.unvested = allocation - claimed - claimable;
    }
}
