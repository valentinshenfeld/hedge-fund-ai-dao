// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.27;

contract TimeLock {
  uint public minDelay;
  address public admin;
  address public pendingAdmin;

  event MinDelayChange(uint oldMinDelay, uint newMinDelay);
  event AdminChange(address oldAdmin, address newAdmin);
  event PendingAdminChange(address oldPendingAdmin, address newPendingAdmin);

  constructor(uint _minDelay) {
    minDelay = _minDelay;
    admin = msg.sender;
    emit AdminChange(address(0), admin);
  }

  function setMinDelay(uint _minDelay) public {
    require(msg.sender == admin, "TimeLock: only admin can set min delay");
    emit MinDelayChange(minDelay, _minDelay);
    minDelay = _minDelay;
  }

  function setPendingAdmin(address _pendingAdmin) public {
    require(msg.sender == admin, "TimeLock: only admin can set pending admin");
    emit PendingAdminChange(pendingAdmin, _pendingAdmin);
    pendingAdmin = _pendingAdmin;
  }

  function acceptAdmin() public {
    require(msg.sender == pendingAdmin, "TimeLock: only pending admin can accept admin");
    emit AdminChange(admin, pendingAdmin);
    admin = pendingAdmin;
    pendingAdmin = address(0);
  }
}
