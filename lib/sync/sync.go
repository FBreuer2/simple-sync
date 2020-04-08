package sync

import (
	"io"
	"os"

	"github.com/FBreuer2/librsync-go"
	"github.com/sirupsen/logrus"

	"golang.org/x/crypto/blake2b"
)

type FileWatcher struct {
	filePath          string
	currentShortState *ShortFileMetadata
	currentFullState  *ExtendedFileMetadata
	changedCallback   func()
}

func NewFileWatcher(path string) (newFileWatcher *FileWatcher, err error) {
	var newWatcher = &FileWatcher{
		filePath: path,
	}

	shortData, err := newWatcher.GetShortFileMetadata()

	if err != nil {
		return nil, err
	}

	newWatcher.currentShortState = shortData
	return newWatcher, nil
}

func (fileWatcher *FileWatcher) ResetCache() {
	fileWatcher.currentShortState = nil
	fileWatcher.currentFullState = nil
}

func (fileWatcher *FileWatcher) GetShortFileMetadata() (metadata *ShortFileMetadata, err error) {
	if fileWatcher.currentShortState != nil {
		return fileWatcher.currentShortState, nil
	}

	inputFile, errInputFile := os.Open(fileWatcher.filePath)

	if errInputFile != nil {
		return nil, errInputFile
	}

	defer inputFile.Close()

	fileInfo, statError := inputFile.Stat()

	if statError != nil {
		return nil, statError
	}

	hasher, err := blake2b.New256(nil)

	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(hasher, inputFile); err != nil {
		return nil, err
	}

	hash := hasher.Sum(nil)

	fileWatcher.currentShortState = &ShortFileMetadata{
		FileSize:    uint64(fileInfo.Size()),
		FileHash:    hash,
		LastChanged: fileInfo.ModTime(),
	}

	return fileWatcher.currentShortState, nil
}

func (fileWatcher *FileWatcher) GetCompleteFileInformation(blockLength uint32, strongChecksumLength uint32) (metadata *ExtendedFileMetadata, err error) {
	if fileWatcher.currentFullState != nil {
		return fileWatcher.currentFullState, nil
	}

	// open the file and check it for more information
	inputFile, errInputFile := os.Open(fileWatcher.filePath)

	if errInputFile != nil {
		return nil, errInputFile
	}

	defer inputFile.Close()

	fileInfo, statError := inputFile.Stat()

	if statError != nil {
		return nil, statError
	}

	signatureFile, err := os.OpenFile(fileWatcher.filePath+".sig", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		logrus.Fatal(err)
	}
	defer signatureFile.Close()

	fileSignatureData, err := librsync.Signature(inputFile, signatureFile, blockLength, strongChecksumLength, librsync.BLAKE2_SIG_MAGIC)

	if err != nil {
		logrus.Fatal(err)
	}

	fileWatcher.currentFullState = &ExtendedFileMetadata{
		FileSize:             uint64(fileInfo.Size()),
		StrongChecksumLength: strongChecksumLength,
		BlockLength:          blockLength,
		BlockAmount:          uint64(len(fileSignatureData.GetStrongChecksums())),
		WeakBlockHashes:      fileSignatureData.GetWeakRollsum(),
		StrongBlockHashes:    fileSignatureData.GetStrongChecksums(),
	}

	return fileWatcher.currentFullState, nil
}
