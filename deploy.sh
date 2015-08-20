#!/bin/bash                                                                                                                                                              
                                                                                                                                                                         
usage() {                                                                                                                                                                
echo "usage: $0 <major> <minor> <release>"                                                                                                                            
}                                                                                                                                                                        
if [ ! $# -eq 3 ]; then                                                                                                                                                  
    usage                                                                                                                                                                 
    exit 1                                                                                                                                                                
fi                                                                                                                                                                       
CWD=$(pwd)                                                                                                                                                               
SVN_ROOT=$RDS_HOME                                                                                                                                                       
MAJOR=$1                                                                                                                                                                 
MINOR=$2                                                                                                                                                                 
RELEASE=$3                                                                                                                                                               
RELEASE_TMP="/tmp/spacex_eggo_${MAJOR}_${MINOR}_${RELEASE}/nice/"                                                                                                              
                                                                                                                                                                          
git pull                                                                                                                                                                 
pushd eggo_spacex                                                                                                                                                        
bash gradlew build --info                                                                                                                                                       
popd                                                                                                                                                                     
                                                                                                                                                                          
rm -rf $RELEASE_TMP                                                                                                                                                      
mkdir -p $RELEASE_TMP                                                                                                                                                    
cp $CWD/predict/prey/ssq.py $RELEASE_TMP/ssq.py                                                                                                                          
cp $CWD/predict/prey/randpicker.py $RELEASE_TMP/randpicker.py                                                                                                            
cp $CWD/predict/prey/logger.conf $RELEASE_TMP/logger.conf                                                                                                                
cp $CWD/eggo_spacex/build/libs/spacex.com-0.1.0.jar $RELEASE_TMP/                                                                
#remove unnecessary pkg                                                                                                                                                  
pushd $RELEASE_TMP/                                                                                                                                                       
tar cvzf $CWD/spacex-${MAJOR}.${MINOR}.${RELEASE}.tar.gz ./..                                                                                                              
popd 
