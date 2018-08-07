npm run build
rm -rf ../../views/index.html
rm -rf ../../static/dist
cp dist/index.html ../../views/index.html
cp -r dist/dist ../../static/dist
cp src/images/avator.jpg ../../static/dist/avator.jpg
cp src/images/login_bg.jpg ../../static/dist/login_bg.jpg
cd ../../ && sh scp.sh
