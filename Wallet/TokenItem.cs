using System;
using System.Globalization;
using System.IO;
using System.Threading.Tasks;
using Cysharp.Threading.Tasks;
using Solana.Unity.Extensions.Models.TokenMint;
using Solana.Unity.Rpc.Models;
using Solana.Unity.SDK.Example;
using Solana.Unity.SDK.Nft;
using Solana.Unity.SDK.Utility;
using TMPro;
using UnityEngine;
using UnityEngine.UI;
using WebSocketSharp;

public class TokenItem : MonoBehaviour
{
    public TextMeshProUGUI pub_txt;
    public TextMeshProUGUI ammount_txt;

    public RawImage logo;
    public Button transferButton;

    public TokenAccount TokenAccount;
    private Nft _nft;

    public TradingPanel tradingPanel;

    private Texture2D _texture;



    //public TextMeshProUGUI title;
    //public TextMeshProUGUI value;
    //public Button Transfer;
    //public Button Receive;


    // Start is called before the first frame update
    void Start()
    {
        logo = GetComponentInChildren<RawImage>();
        transferButton.onClick.AddListener(TransferAccount);
    }

    // Update is called once per frame
    void Update()
    {

    }

    public async void InitPanel(TokenAccount tokenAccount, Nft nftData = null)
    {
        TokenAccount = tokenAccount;
        if (nftData != null && ulong.Parse(tokenAccount.Account.Data.Parsed.Info.TokenAmount.Amount) == 1)
        {
            await UniTask.SwitchToMainThread();
            _nft = nftData;
            ammount_txt.text = "";
            pub_txt.text = nftData.metaplexData?.data?.offchainData?.name;

            if (logo != null)
            {
                logo.texture = nftData.metaplexData?.nftImage?.file;
            }
        }
        else
        {
            ammount_txt.text =
    tokenAccount.Account.Data.Parsed.Info.TokenAmount.AmountDecimal.ToString(CultureInfo
        .CurrentCulture);
            pub_txt.text = nftData?.metaplexData?.data?.offchainData?.name ?? tokenAccount.Account.Data.Parsed.Info.Mint;
            if (pub_txt.text.ToLower().Contains("igames"))
            {
                PlayerInfo.instance.IGSTokenAccount = tokenAccount;
            }
            if (nftData?.metaplexData?.data?.offchainData?.symbol != null)
            {
                pub_txt.text += $" ({nftData?.metaplexData?.data?.offchainData?.symbol})";
            }

            if (nftData?.metaplexData?.data?.offchainData?.default_image != null)
            {
                await LoadAndCacheTokenLogo(nftData.metaplexData?.data?.offchainData?.default_image, tokenAccount.Account.Data.Parsed.Info.Mint);
            }
            else
            {
                var tokenMintResolver = await WalletScreen.GetTokenMintResolver();
                TokenDef tokenDef = tokenMintResolver.Resolve(tokenAccount.Account.Data.Parsed.Info.Mint);
                if (tokenDef.TokenName.IsNullOrEmpty() || tokenDef.Symbol.IsNullOrEmpty()) return;
                pub_txt.text = $"{tokenDef.TokenName} ({tokenDef.Symbol})";
                await LoadAndCacheTokenLogo(tokenDef.TokenLogoUrl, tokenDef.TokenMint);
            }

        }

    }

    private async Task LoadAndCacheTokenLogo(string logoUrl, string tokenMint)
    {
        if (logoUrl.IsNullOrEmpty() || tokenMint.IsNullOrEmpty() || logo is null) return;
        var texture = await FileLoader.LoadFile<Texture2D>(logoUrl);
        _texture = FileLoader.Resize(texture, 75, 75);
        FileLoader.SaveToPersistentDataPath(Path.Combine(Application.persistentDataPath, $"{tokenMint}.png"), _texture);
        logo.texture = _texture;
    }

    public void TransferAccount()
    {
        if (_nft != null)
        {
            TradingPanel panel = UnityEngine.Object.Instantiate(tradingPanel, transform.root).GetComponent<TradingPanel>();

            panel.Show(Common.TradingType.Transfer,null,_nft);
        }
        else if(TokenAccount!=null)
        {
            TradingPanel panel = UnityEngine.Object.Instantiate(tradingPanel, transform.root).GetComponent<TradingPanel>();

            panel.Show(Common.TradingType.Transfer, TokenAccount);
        }
        else
        {
            TradingPanel panel = UnityEngine.Object.Instantiate(tradingPanel, transform.root).GetComponent<TradingPanel>();

            panel.Show(Common.TradingType.Transfer);
        }
    }

}
