cd backend
xfce4-terminal -e "go run ." --hold --geometry 66x16-0-0
cd ..

cd frontend
xfce4-terminal -e "npm run serve" --hold --geometry 66x16-0+0

firefox --new-tab http://127.0.0.1:8080