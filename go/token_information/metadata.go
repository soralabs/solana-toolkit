package token_information

import (
	"context"
	"encoding/binary"

	"github.com/gagliardetto/solana-go"
)

const METADATA_PROGRAM_ID = "metaqbxxUerdq28cj1RbAWkYQm3ybzjb6a8bt518x1s"

// deriveMetadataAddress derives the metadata address for a given token address
func deriveMetadataAddress(tokenAddress solana.PublicKey) (solana.PublicKey, error) {
	metadataProgramID, err := solana.PublicKeyFromBase58(METADATA_PROGRAM_ID)
	if err != nil {
		return solana.PublicKey{}, err
	}

	seeds := [][]byte{
		[]byte("metadata"),
		metadataProgramID.Bytes(),
		tokenAddress.Bytes(),
	}

	addr, _, err := solana.FindProgramAddress(seeds, metadataProgramID)
	return addr, err
}

// decodeMetadata decodes the metadata from the given data
func decodeMetadata(data []byte) (Metadata, error) {
	metadata := &Metadata{}
	metadata.Key = data[0]
	metadata.UpdateAuthority = solana.PublicKeyFromBytes(data[1:33])
	metadata.Mint = solana.PublicKeyFromBytes(data[33:65])

	readUint32 := func(offset int) (uint32, int) {
		return binary.LittleEndian.Uint32(data[offset : offset+4]), offset + 4
	}

	nameLen, offset := readUint32(65)
	symbolLen, offset := readUint32(offset + int(nameLen))
	uriLen, offset := readUint32(offset + int(symbolLen))

	metadata.Data.Name = string(data[69 : 69+nameLen])
	metadata.Data.Symbol = string(data[73+nameLen : 73+nameLen+symbolLen])
	metadata.Data.Uri = string(data[77+nameLen+symbolLen : 77+nameLen+symbolLen+uriLen])

	offset = int(77 + nameLen + symbolLen + uriLen)
	metadata.Data.SellerFeeBasisPoints = binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2

	if len(data) > int(offset) {
		metadata.Data.CreatorCount = data[offset]
		offset++

		// Preallocate slice for creators
		metadata.Data.Creators = make([]MetadataCreator, metadata.Data.CreatorCount)
		for i := uint8(0); i < metadata.Data.CreatorCount; i++ {
			creator := &metadata.Data.Creators[i]
			creator.Address = solana.PublicKeyFromBytes(data[offset : offset+32])
			offset += 32
			creator.Verified = data[offset] != 0
			offset++
			creator.Share = data[offset]
			offset++
		}
	}

	if len(data) > int(offset) {
		metadata.PrimarySaleHappened = data[offset] != 0
		offset++
	}

	if len(data) > int(offset) {
		metadata.IsMutable = data[offset] != 0
	}

	return *metadata, nil
}

// getMetadata fetches the metadata on chain for a given token address
func (t *TokenInformationTool) getMetadata(ctx context.Context, tokenAddress solana.PublicKey) (Metadata, error) {
	metadataAddress, err := deriveMetadataAddress(tokenAddress)
	if err != nil {
		return Metadata{}, err
	}

	metadata, err := t.rpcClient.GetAccountInfo(ctx, metadataAddress)
	if err != nil {
		return Metadata{}, err
	}

	return decodeMetadata(metadata.Value.Data.GetBinary())
}
