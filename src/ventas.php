
<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.6.1/jquery.min.js"></script>
<?php include_once "includes/header.php";
require("../conexion.php");

$id_user = $_SESSION['idUser'];
$permiso = "nueva_venta";
$sql = mysqli_query($conexion, "SELECT p.*, d.* FROM permisos p INNER JOIN detalle_permisos d ON p.id = d.id_permiso WHERE d.id_usuario = $id_user AND p.nombre = '$permiso'");
$existe = mysqli_fetch_all($sql);
if (empty($existe) && $id_user != 1) {
    header("Location: permisos.php");
}
?>
<div class="row">
    <div class="col-lg-12">
        <div class="form-group">
            <h4 class="text-center">Datos del Cliente</h4>
        </div>
        <div class="card">
            <div class="card-body">
                <form method="post">
                    <div class="row">
                        <div class="col-lg-4">
                            <div>
                                <input type="hidden" id="idcliente" value="1" name="idcliente" required>
                                <label>Nombre</label>
                                <input type="text" name="nom_cliente" id="nom_cliente" class="form-control" placeholder="Ingrese nombre del cliente" required>
                            </div>
                        </div>
                        <div class="col-lg-4">
                            <div class="form-group">
                                <label>Teléfono</label>
                                <input type="number" name="tel_cliente" id="tel_cliente" class="form-control" disabled required>
                            </div>
                        </div>
                        <div class="col-lg-4">
                            <div class="form-group">
                                <label>Dirreción</label>
                                <input type="text" name="dir_cliente" id="dir_cliente" class="form-control" disabled required>
                            </div>
                        </div>
                    </div>
                </form>
            </div>
        </div>
        <div class="card">
            <div class="card-header bg-primary text-white text-center">
                Datos Venta
            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-lg-6">
                        <div class="form-group">
                            <i class="fas fa-user"></i><p style="font-size: 16px; text-transform: uppercase; color: red;"><?php echo $_SESSION['nombre']; ?></p>
                            <label> <b>Graduaciones:</b></label>
                            <p style="font-size: 16px; text-transform: uppercase; color: black;">
                            <table class="table table-borderless" id="tablaGracuadiones">
                            <form id="graduaciones">
                                <div id="okgrad">
                                <tr>
                                    <td><b>Graduacion Lejos </b></td>
                                </tr>
                                <tr>
                                    <td><div class="col-xs-3"><b>Ojo D: </b><input  id="ojoDl1" name="ojoDl1" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoDl2" name="ojoDl2" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoDl3" name="ojoDl3" type="text" size="4"></div></td>
                                </tr>
                                <tr>
                                    <td><b>Ojo I: &nbsp</b><input id="ojoIl1" name="ojoIl1" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoIl2" name="ojoIl2" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoIl3" name="ojoIl3" type="text" size="4"></td>
                                </tr>  
                                <tr>
                                    <td>&nbsp<b>ADD: &nbsp</b><input id="add" name="add" type="text" size="4">&nbsp&nbsp&nbsp<b>Obs: &nbsp</b><input id="obs" name="obs" type="text" size="10"></td>
                                </tr>  
                                <tr>
                                    
                                    <td><b>Graduacion Cerca </b></td>
                                </tr>
                                
                                <tr>
                                    <td><b>Ojo D: </b><input name="ojoD1" id="ojoD1" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoD2" name="ojoD2" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoD3" name="ojoD3" type="text" size="4"></td>
                                </tr>
                                <tr>
                                    <td><b>Ojo I: &nbsp</b><input id="ojoI1" name="ojoI1" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoI2" name="ojoI2" type="text" size="4">&nbsp&nbsp&nbsp<input id="ojoI3" name="ojoI3" type="text" size="4"></td>
                                </tr>
                                
                                <tr>
                                    <td><input class="btn btn-primary" name="grad" id="grad" type="button" value="Agregar"></td>
                                </tr>
                                </div>
                            </form>
                        </table>    
                        </p>
                        </div>
                    </div>
                    <div class="col-lg-6">
                        <div class="card">
                            <div class="card-header">
                                Buscar Producto
                            </div>
                            <div class="card-body">
                                <div class="form-group">
                                    <input id="producto" class="form-control" type="text" name="producto" placeholder="Ingresa el código o nombre">
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

            </div>
        </div>
        <div class="table-responsive">
            <table class="table table-hover" id="tblDetalle">
                <thead class="thead-dark">
                    <tr>
                        <th>Id</th>
                        <th>Descripción</th>
                        <th>Cantidad</th>
                        <th>Precio</th>
                        <th>Precio Total</th>
                        <th>Accion</th>
                    </tr>
                </thead>
                <tbody id="detalle_venta">

                </tbody>
                <tfoot>
                    <tr class="font-weight-bold">
                        <td colspan=3>Total a Pagar: </td>
                        <td></td>
                    </tr>    
                </tfoot>
            </table>
            <table class="table table-hover">
                <form method="POST" id="form_descuento">
                    <div id="div_descuento">
                    <tr class="font-weight-bold">
                        <td colspan=3>Abona: </td>
                        <td colspan=3><input type="text" size="3"> </td>
                    </tr>
                    <tr class="font-weight-bold">
                        <td colspan=3>Descuento: </td>
                        <td colspan=3><select id="porc" name="porc">
                                        <option value="1">Sin descuento</option>
                                        <option value="0.95">5%</option>
                                        <option value="0.9">10%</option>
                                        <option value="0.85">15%</option>
                                        <option value="0.80">20%</option>
                                        <option value="0.75">25%</option>
                                        <option value="0.70">30%</option>
                                        <option value="0.65">35%</option>
                                        <option value="0.60">40%</option>
                                        <option value="0.55">45%</option>
                                        <option value="0.50">50%</option>
                                        <option value="0.45">55%</option>
                                        <option value="0.40">60%</option>
                                        </select></td>
                    </tr>
                    <tr class="font-weight-bold">
                        <td colspan=3>Resta: </td>
                        <td colspan=3><input type="text" size="3" disabled></td>
                        </tr>
                </tr>
                </div>
                </form>
            </table>
        </div>
    </div>
    <div class="col-md-6">
    <a href="#" class="btn btn-primary" id="btn_generar"><i class="fas fa-save"></i> Generar Venta</a>
    <input type="button" class="btn btn-primary" value="Aplicar Descuento" name="btn_descuento" id="btn_descuento" onclick=descuento()></input>
    </div>

</div>



<?php include_once "includes/footer.php"; ?>



