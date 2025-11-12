// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.27;

import "forge-std/Script.sol";
import "forge-std/console.sol";

interface IMATPFactory {
    function createMATPs(
        address[] memory _beneficiaries,
        uint256[] memory _allocations,
        uint256[] memory _milestoneIds
    ) external returns (address[] memory);
}

contract CreateMATPScript is Script{
    function run() external {
        // Load your private key from environment variable
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        // Start broadcasting transactions
        vm.startBroadcast(deployerPrivateKey);

        // ATPFactory address on Sepolia
        address factoryAddress = 0x942cDD7C5620FA7476e0359b78c69AB654D789f4;
        IMATPFactory factory = IMATPFactory(factoryAddress);

        // Declare dynamic memory arrays
        address ;
        uint256 ;
        uint256 ;

        // Fill in the arrays
        beneficiaries[0] = 0x5A9dadc51D863Bc7124C9eD3C4894cA41215b86C;
        allocations[0] = 1 ether; // 1 ETH
        milestoneIds[0] = 1;

        // Call the createMATPs function
        address[] memory matps = factory.createMATPs(
            beneficiaries,
            allocations,
            milestoneIds
        );

        // Log the created MATP addresses
        for (uint256 i = 0; i < matps.length; i++) {
            console.log("MATP created:", matps[i]);
        }

        vm.stopBroadcast();
    }
}
