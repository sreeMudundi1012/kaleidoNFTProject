// SPDX-License-Identifier: MIT
pragma solidity >=0.4.20 <0.9.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";

contract KaleidoNFT is ERC721URIStorage, Ownable{
    using Counters for Counters.Counter;
    Counters.Counter private newTokenID;

    constructor() ERC721("Kaleido NFT", "KLD"){
    }

    function mintNFT(string memory tokenURI) public onlyOwner returns (uint256){
        newTokenID.increment();
        uint256 newID = newTokenID.current();
        _mint(msg.sender, newID);
        _setTokenURI(newID, tokenURI);
        return newID;
    }

    function burnNFT(uint256 tokenID) public {
        _burn(tokenID);
    }

    function transferNFT(uint256 tokenID, address from, address to )public onlyOwner{
        transferFrom(from, to, tokenID);
    }
}